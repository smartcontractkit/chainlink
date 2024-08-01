package serialization

import (
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr2/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"google.golang.org/protobuf/proto"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

// Serialize encodes a protocol.Message into a binary payload
func Serialize(m protocol.Message) (b []byte, pbm *MessageWrapper, err error) {
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
func Deserialize(b []byte) (protocol.Message, *MessageWrapper, error) {
	pbm := &MessageWrapper{}
	err := proto.Unmarshal(b, pbm)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal protobuf: %w", err)
	}
	m, err := messageWrapperFromProtoMessage(pbm)
	if err != nil {
		return nil, nil, fmt.Errorf("could not translate protobuf to protocol.Message: %w", err)
	}
	return m, pbm, nil
}

//
// *toProtoMessage
//

func toProtoMessage(m protocol.Message) (*MessageWrapper, error) {
	msgWrapper := MessageWrapper{}
	switch v := m.(type) {
	case protocol.MessageNewEpoch:
		pm := &MessageNewEpoch{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
		}
		msgWrapper.Msg = &MessageWrapper_MessageNewEpoch{pm}
	case protocol.MessageObserveReq:
		pm := &MessageObserveReq{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			uint32(v.Round),
			v.Query,
		}
		msgWrapper.Msg = &MessageWrapper_MessageObserveReq{pm}
	case protocol.MessageObserve:
		pm := &MessageObserve{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			uint32(v.Round),
			signedObservationToProtoMessage(v.SignedObservation),
		}
		msgWrapper.Msg = &MessageWrapper_MessageObserve{pm}
	case protocol.MessageReportReq:
		pbasos := make([]*AttributedSignedObservation, 0, len(v.AttributedSignedObservations))
		for _, aso := range v.AttributedSignedObservations {
			pbasos = append(pbasos, attributedSignedObservationToProtoMessage(aso))
		}
		pm := &MessageReportReq{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			uint32(v.Round),
			v.Query,
			pbasos,
		}
		msgWrapper.Msg = &MessageWrapper_MessageReportReq{pm}
	case protocol.MessageReport:
		pm := &MessageReport{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			uint32(v.Round),
			attestedReportOneToProtoMessage(v.AttestedReport),
		}
		msgWrapper.Msg = &MessageWrapper_MessageReport{pm}
	case protocol.MessageFinal:
		msgWrapper.Msg = &MessageWrapper_MessageFinal{finalToProtoMessage(v)}
	case protocol.MessageFinalEcho:
		msgWrapper.Msg = &MessageWrapper_MessageFinalEcho{
			&MessageFinalEcho{
				// zero-initialize protobuf built-ins
				protoimpl.MessageState{},
				0,
				nil,
				// fields
				finalToProtoMessage(v.MessageFinal),
			},
		}
	default:
		return nil, fmt.Errorf("unable to serialize message of type %T", m)

	}
	return &msgWrapper, nil
}

func signedObservationToProtoMessage(o protocol.SignedObservation) *SignedObservation {
	return &SignedObservation{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		o.Observation,
		o.Signature,
	}
}

func attributedSignedObservationToProtoMessage(aso protocol.AttributedSignedObservation) *AttributedSignedObservation {
	return &AttributedSignedObservation{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		signedObservationToProtoMessage(aso.SignedObservation),
		uint32(aso.Observer),
	}
}

func attestedReportOneToProtoMessage(aro protocol.AttestedReportOne) *AttestedReportOne {
	return &AttestedReportOne{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		aro.Skip,
		aro.Report,
		aro.Signature,
	}
}

func attestedReportManyToProtoMessage(arm protocol.AttestedReportMany) *AttestedReportMany {
	pbass := make([]*AttributedSignature, 0, len(arm.AttributedSignatures))
	for _, as := range arm.AttributedSignatures {
		pbass = append(pbass, &AttributedSignature{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			as.Signature,
			uint32(as.Signer),
		})
	}
	return &AttestedReportMany{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		arm.Report,
		pbass,
	}
}

func finalToProtoMessage(v protocol.MessageFinal) *MessageFinal {
	return &MessageFinal{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		uint64(v.Epoch),
		uint32(v.Round),
		v.H[:],
		attestedReportManyToProtoMessage(v.AttestedReport),
	}
}

//
// *fromProtoMessage
//

func messageWrapperFromProtoMessage(wrapper *MessageWrapper) (protocol.Message, error) {
	switch msg := wrapper.Msg.(type) {
	case *MessageWrapper_MessageNewEpoch:
		return messageNewEpochFromProtoMessage(wrapper.GetMessageNewEpoch())
	case *MessageWrapper_MessageObserveReq:
		return messageObserveReqFromProtoMessage(wrapper.GetMessageObserveReq())
	case *MessageWrapper_MessageObserve:
		return messageObserveFromProtoMessage(wrapper.GetMessageObserve())
	case *MessageWrapper_MessageReportReq:
		return messageReportReqFromProtoMessage(wrapper.GetMessageReportReq())
	case *MessageWrapper_MessageReport:
		return messageReportFromProtoMessage(wrapper.GetMessageReport())
	case *MessageWrapper_MessageFinal:
		return messageFinalFromProtoMessage(wrapper.GetMessageFinal())
	case *MessageWrapper_MessageFinalEcho:
		return messageFinalEchoFromProtoMessage(wrapper.GetMessageFinalEcho())
	default:
		return nil, fmt.Errorf("unrecognized Msg type %T", msg)
	}
}

func messageNewEpochFromProtoMessage(m *MessageNewEpoch) (protocol.MessageNewEpoch, error) {
	if m == nil {
		return protocol.MessageNewEpoch{}, fmt.Errorf("unable to extract a MessageNewEpoch value")
	}
	return protocol.MessageNewEpoch{
		uint32(m.Epoch),
	}, nil
}

func messageObserveReqFromProtoMessage(m *MessageObserveReq) (protocol.MessageObserveReq, error) {
	if m == nil {
		return protocol.MessageObserveReq{}, fmt.Errorf("unable to extract a MessageObserveReq value")
	}
	return protocol.MessageObserveReq{
		uint32(m.Epoch),
		uint8(m.Round),
		m.Query,
	}, nil
}

func messageObserveFromProtoMessage(m *MessageObserve) (protocol.MessageObserve, error) {
	if m == nil {
		return protocol.MessageObserve{}, fmt.Errorf("unable to extract a MessageObserve value")
	}
	so, err := signedObservationFromProtoMessage(m.SignedObservation)
	if err != nil {
		return protocol.MessageObserve{}, err
	}
	return protocol.MessageObserve{
		uint32(m.Epoch),
		uint8(m.Round),
		so,
	}, nil
}

func messageReportReqFromProtoMessage(m *MessageReportReq) (protocol.MessageReportReq, error) {
	if m == nil {
		return protocol.MessageReportReq{}, fmt.Errorf("unable to extract a MessageReportReq value")
	}
	asos, err := attributedSignedObservationsFromProtoMessage(m.AttributedSignedObservations)
	if err != nil {
		return protocol.MessageReportReq{}, err
	}
	return protocol.MessageReportReq{
		uint32(m.Epoch),
		uint8(m.Round),
		m.Query,
		asos,
	}, nil
}

func attestedReportOneFromProtoMessage(m *AttestedReportOne) (protocol.AttestedReportOne, error) {
	if m == nil {
		return protocol.AttestedReportOne{}, fmt.Errorf("unable to extract a AttestedReportOne value")
	}

	return protocol.AttestedReportOne{
		m.Skip,
		m.Report,
		m.Signature,
	}, nil
}

func messageReportFromProtoMessage(m *MessageReport) (protocol.MessageReport, error) {
	if m == nil {
		return protocol.MessageReport{}, fmt.Errorf("unable to extract a MessageReport value")
	}

	report, err := attestedReportOneFromProtoMessage(m.AttestedReport)
	if err != nil {
		return protocol.MessageReport{}, err
	}

	return protocol.MessageReport{
		uint32(m.Epoch),
		uint8(m.Round),
		report,
	}, nil
}

func attestedReportManyFromProtoMessage(m *AttestedReportMany) (protocol.AttestedReportMany, error) {
	if m == nil {
		return protocol.AttestedReportMany{}, fmt.Errorf("unable to extract a AttestedReportMany value")
	}

	ass := make([]types.AttributedOnchainSignature, 0, len(m.AttributedSignatures))
	for i, as := range m.AttributedSignatures {
		if as == nil {
			return protocol.AttestedReportMany{}, fmt.Errorf("unable to extract a AttestedReportMany value because AttributedSignatures[%v] is nil", i)
		}
		ass = append(ass, types.AttributedOnchainSignature{
			as.Signature,
			commontypes.OracleID(as.Signer),
		})
	}

	return protocol.AttestedReportMany{
		m.Report,
		ass,
	}, nil
}

func messageFinalFromProtoMessage(m *MessageFinal) (protocol.MessageFinal, error) {
	if m == nil {
		return protocol.MessageFinal{}, fmt.Errorf("unable to extract a MessageFinal value")
	}
	report, err := attestedReportManyFromProtoMessage(m.AttestedReport)
	if err != nil {
		return protocol.MessageFinal{}, err
	}
	var h [32]byte
	if len(m.H) != len(h) {
		return protocol.MessageFinal{}, fmt.Errorf("wrong length for MessageFinal.H. got %v but wanted %v", len(m.H), len(h))
	}
	copy(h[:], m.H)
	return protocol.MessageFinal{
		uint32(m.Epoch),
		uint8(m.Round),
		h,
		report,
	}, nil
}

func messageFinalEchoFromProtoMessage(m *MessageFinalEcho) (protocol.MessageFinalEcho, error) {
	if m == nil {
		return protocol.MessageFinalEcho{}, fmt.Errorf("unable to extract a MessageFinalEcho value")
	}
	final, err := messageFinalFromProtoMessage(m.Final)
	if err != nil {
		return protocol.MessageFinalEcho{}, err
	}
	return protocol.MessageFinalEcho{final}, nil
}

func attributedSignedObservationsFromProtoMessage(pbasos []*AttributedSignedObservation) ([]protocol.AttributedSignedObservation, error) {
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

func attributedSignedObservationFromProtoMessage(m *AttributedSignedObservation) (protocol.AttributedSignedObservation, error) {
	if m == nil {
		return protocol.AttributedSignedObservation{}, fmt.Errorf("unable to extract an AttributedSignedObservation value")
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

func signedObservationFromProtoMessage(m *SignedObservation) (protocol.SignedObservation, error) {
	if m == nil {
		return protocol.SignedObservation{}, fmt.Errorf("unable to extract an SignedObservation value")
	}

	return protocol.SignedObservation{
		m.Observation,
		m.Signature,
	}, nil
}
