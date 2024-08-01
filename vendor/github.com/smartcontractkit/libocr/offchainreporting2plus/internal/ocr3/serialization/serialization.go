package serialization

import (
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol"

	"google.golang.org/protobuf/proto"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

// Serialize encodes a protocol.Message into a binary payload
func Serialize[RI any](m protocol.Message[RI]) (b []byte, pbm *MessageWrapper, err error) {
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
func Deserialize[RI any](b []byte) (protocol.Message[RI], *MessageWrapper, error) {
	pbm := &MessageWrapper{}
	err := proto.Unmarshal(b, pbm)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal protobuf: %w", err)
	}
	m, err := messageWrapperFromProtoMessage[RI](pbm)
	if err != nil {
		return nil, nil, fmt.Errorf("could not translate protobuf to protocol.Message: %w", err)
	}
	return m, pbm, nil
}

//
// *toProtoMessage
//

func toProtoMessage[RI any](m protocol.Message[RI]) (*MessageWrapper, error) {
	msgWrapper := MessageWrapper{}
	switch v := m.(type) {
	case protocol.MessageNewEpochWish[RI]:
		pm := &MessageNewEpochWish{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
		}
		msgWrapper.Msg = &MessageWrapper_MessageNewEpochWish{pm}
	case protocol.MessageEpochStartRequest[RI]:
		pm := &MessageEpochStartRequest{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			CertifiedPrepareOrCommitToProtoMessage(v.HighestCertified),
			signedHighestCertifiedTimestampToProtoMessage(v.SignedHighestCertifiedTimestamp),
		}
		msgWrapper.Msg = &MessageWrapper_MessageEpochStartRequest{pm}
	case protocol.MessageEpochStart[RI]:
		pm := &MessageEpochStart{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			epochStartProofToProtoMessage(v.EpochStartProof),
		}
		msgWrapper.Msg = &MessageWrapper_MessageEpochStart{pm}
	case protocol.MessageRoundStart[RI]:
		pm := &MessageRoundStart{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			v.Query,
		}
		msgWrapper.Msg = &MessageWrapper_MessageRoundStart{pm}
	case protocol.MessageObservation[RI]:
		pm := &MessageObservation{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			signedObservationToProtoMessage(v.SignedObservation),
		}
		msgWrapper.Msg = &MessageWrapper_MessageObservation{pm}

	case protocol.MessageProposal[RI]:
		pbasos := make([]*AttributedSignedObservation, 0, len(v.AttributedSignedObservations))
		for _, aso := range v.AttributedSignedObservations {
			pbasos = append(pbasos, attributedSignedObservationToProtoMessage(aso))
		}
		pm := &MessageProposal{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			pbasos,
		}
		msgWrapper.Msg = &MessageWrapper_MessageProposal{pm}
	case protocol.MessagePrepare[RI]:
		pm := &MessagePrepare{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			v.Signature,
		}
		msgWrapper.Msg = &MessageWrapper_MessagePrepare{pm}
	case protocol.MessageCommit[RI]:
		pm := &MessageCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			v.Signature,
		}
		msgWrapper.Msg = &MessageWrapper_MessageCommit{pm}
	case protocol.MessageReportSignatures[RI]:
		pm := &MessageReportSignatures{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.SeqNr,
			v.ReportSignatures,
		}
		msgWrapper.Msg = &MessageWrapper_MessageReportSignatures{pm}
	case protocol.MessageCertifiedCommitRequest[RI]:
		pm := &MessageCertifiedCommitRequest{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.SeqNr,
		}
		msgWrapper.Msg = &MessageWrapper_MessageCertifiedCommitRequest{pm}
	case protocol.MessageCertifiedCommit[RI]:
		pm := &MessageCertifiedCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			CertifiedCommitToProtoMessage(v.CertifiedCommit),
		}
		msgWrapper.Msg = &MessageWrapper_MessageCertifiedCommit{pm}

	default:
		return nil, fmt.Errorf("unable to serialize message of type %T", m)

	}
	return &msgWrapper, nil
}

func CertifiedPrepareOrCommitToProtoMessage(cpoc protocol.CertifiedPrepareOrCommit) *CertifiedPrepareOrCommit {
	switch v := cpoc.(type) {
	case *protocol.CertifiedPrepare:
		prepareQuorumCertificate := make([]*AttributedPrepareSignature, 0, len(v.PrepareQuorumCertificate))
		for _, aps := range v.PrepareQuorumCertificate {
			prepareQuorumCertificate = append(prepareQuorumCertificate, &AttributedPrepareSignature{
				// zero-initialize protobuf built-ins
				protoimpl.MessageState{},
				0,
				nil,
				// fields
				aps.Signature,
				uint32(aps.Signer),
			})
		}
		return &CertifiedPrepareOrCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			&CertifiedPrepareOrCommit_Prepare{&CertifiedPrepare{
				// zero-initialize protobuf built-ins
				protoimpl.MessageState{},
				0,
				nil,
				// fields
				uint64(v.PrepareEpoch),
				v.SeqNr,
				v.OutcomeInputsDigest[:],
				v.Outcome,
				prepareQuorumCertificate,
			}},
		}
	case *protocol.CertifiedCommit:
		return &CertifiedPrepareOrCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			&CertifiedPrepareOrCommit_Commit{CertifiedCommitToProtoMessage(*v)},
		}
	default:
		// It's safe to crash here since the "protocol.*" versions of these values
		// come from the trusted, local environment.
		panic("unrecognized protocol.CertifiedPrepareOrCommit implementation")
	}
}

func CertifiedCommitToProtoMessage(cpocc protocol.CertifiedCommit) *CertifiedCommit {
	commitQuorumCertificate := make([]*AttributedCommitSignature, 0, len(cpocc.CommitQuorumCertificate))
	for _, aps := range cpocc.CommitQuorumCertificate {
		commitQuorumCertificate = append(commitQuorumCertificate, &AttributedCommitSignature{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			aps.Signature,
			uint32(aps.Signer),
		})
	}
	return &CertifiedCommit{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		uint64(cpocc.CommitEpoch),
		cpocc.SeqNr,
		cpocc.Outcome,
		commitQuorumCertificate,
	}
}

func attributedSignedHighestCertifiedTimestampToProtoMessage(ashct protocol.AttributedSignedHighestCertifiedTimestamp) *AttributedSignedHighestCertifiedTimestamp {
	return &AttributedSignedHighestCertifiedTimestamp{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		signedHighestCertifiedTimestampToProtoMessage(ashct.SignedHighestCertifiedTimestamp),
		uint32(ashct.Signer),
	}
}

func signedHighestCertifiedTimestampToProtoMessage(shct protocol.SignedHighestCertifiedTimestamp) *SignedHighestCertifiedTimestamp {
	return &SignedHighestCertifiedTimestamp{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		highestCertifiedTimestampToProtoMessage(shct.HighestCertifiedTimestamp),
		shct.Signature,
	}
}

func highestCertifiedTimestampToProtoMessage(hct protocol.HighestCertifiedTimestamp) *HighestCertifiedTimestamp {
	return &HighestCertifiedTimestamp{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		hct.SeqNr,
		hct.CommittedElsePrepared,
	}
}

func epochStartProofToProtoMessage(srqc protocol.EpochStartProof) *EpochStartProof {
	highestCertifiedProof := make([]*AttributedSignedHighestCertifiedTimestamp, 0, len(srqc.HighestCertifiedProof))
	for _, ashct := range srqc.HighestCertifiedProof {
		highestCertifiedProof = append(highestCertifiedProof, attributedSignedHighestCertifiedTimestampToProtoMessage(ashct))
	}
	return &EpochStartProof{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		CertifiedPrepareOrCommitToProtoMessage(srqc.HighestCertified),
		highestCertifiedProof,
	}
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

func PacemakerStateToProtoMessage(ps protocol.PacemakerState) *PacemakerState {
	return &PacemakerState{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		ps.Epoch,
		ps.HighestSentNewEpochWish,
	}
}

//
// *fromProtoMessage
//

func messageWrapperFromProtoMessage[RI any](wrapper *MessageWrapper) (protocol.Message[RI], error) {
	switch msg := wrapper.Msg.(type) {
	case *MessageWrapper_MessageNewEpochWish:
		return messageNewEpochWishFromProtoMessage[RI](wrapper.GetMessageNewEpochWish())
	case *MessageWrapper_MessageEpochStartRequest:
		return messageEpochStartRequestFromProtoMessage[RI](wrapper.GetMessageEpochStartRequest())
	case *MessageWrapper_MessageEpochStart:
		return messageEpochStartFromProtoMessage[RI](wrapper.GetMessageEpochStart())
	case *MessageWrapper_MessageRoundStart:
		return messageRoundStartFromProtoMessage[RI](wrapper.GetMessageRoundStart())
	case *MessageWrapper_MessageObservation:
		return messageObservationFromProtoMessage[RI](wrapper.GetMessageObservation())
	case *MessageWrapper_MessageProposal:
		return messageProposalFromProtoMessage[RI](wrapper.GetMessageProposal())
	case *MessageWrapper_MessagePrepare:
		return messagePrepareFromProtoMessage[RI](wrapper.GetMessagePrepare())
	case *MessageWrapper_MessageCommit:
		return messageCommitFromProtoMessage[RI](wrapper.GetMessageCommit())
	case *MessageWrapper_MessageReportSignatures:
		return messageReportSignaturesFromProtoMessage[RI](wrapper.GetMessageReportSignatures())
	case *MessageWrapper_MessageCertifiedCommitRequest:
		return messageCertifiedCommitRequestFromProtoMessage[RI](wrapper.GetMessageCertifiedCommitRequest())
	case *MessageWrapper_MessageCertifiedCommit:
		return messageCertifiedCommitFromProtoMessage[RI](wrapper.GetMessageCertifiedCommit())
	default:
		return nil, fmt.Errorf("unrecognized Msg type %T", msg)
	}
}

func messageNewEpochWishFromProtoMessage[RI any](m *MessageNewEpochWish) (protocol.MessageNewEpochWish[RI], error) {
	if m == nil {
		return protocol.MessageNewEpochWish[RI]{}, fmt.Errorf("unable to extract a MessageNewEpochWish value")
	}
	return protocol.MessageNewEpochWish[RI]{
		m.Epoch,
	}, nil
}

func messageEpochStartRequestFromProtoMessage[RI any](m *MessageEpochStartRequest) (protocol.MessageEpochStartRequest[RI], error) {
	if m == nil {
		return protocol.MessageEpochStartRequest[RI]{}, fmt.Errorf("unable to extract a MessageEpochStartRequest value")
	}
	hc, err := CertifiedPrepareOrCommitFromProtoMessage(m.HighestCertified)
	if err != nil {
		return protocol.MessageEpochStartRequest[RI]{}, err
	}
	shct, err := signedHighestCertifiedTimestampFromProtoMessage(m.SignedHighestCertifiedTimestamp)
	if err != nil {
		return protocol.MessageEpochStartRequest[RI]{}, err
	}
	return protocol.MessageEpochStartRequest[RI]{
		m.Epoch,
		hc,
		shct,
	}, nil
}

func messageEpochStartFromProtoMessage[RI any](m *MessageEpochStart) (protocol.MessageEpochStart[RI], error) {
	if m == nil {
		return protocol.MessageEpochStart[RI]{}, fmt.Errorf("unable to extract a MessageEpochStart value")
	}
	srqc, err := epochStartProofFromProtoMessage(m.EpochStartProof)
	if err != nil {
		return protocol.MessageEpochStart[RI]{}, err
	}
	return protocol.MessageEpochStart[RI]{
		m.Epoch,
		srqc,
	}, nil
}

func messageProposalFromProtoMessage[RI any](m *MessageProposal) (protocol.MessageProposal[RI], error) {
	if m == nil {
		return protocol.MessageProposal[RI]{}, fmt.Errorf("unable to extract a MessageProposal value")
	}
	asos, err := attributedSignedObservationsFromProtoMessage(m.AttributedSignedObservations)
	if err != nil {
		return protocol.MessageProposal[RI]{}, err
	}
	return protocol.MessageProposal[RI]{
		m.Epoch,
		m.SeqNr,
		asos,
	}, nil
}

func messagePrepareFromProtoMessage[RI any](m *MessagePrepare) (protocol.MessagePrepare[RI], error) {
	if m == nil {
		return protocol.MessagePrepare[RI]{}, fmt.Errorf("unable to extract a MessagePrepare value")
	}
	return protocol.MessagePrepare[RI]{
		m.Epoch,
		m.SeqNr,
		m.Signature,
	}, nil
}

func messageCommitFromProtoMessage[RI any](m *MessageCommit) (protocol.MessageCommit[RI], error) {
	if m == nil {
		return protocol.MessageCommit[RI]{}, fmt.Errorf("unable to extract a MessageCommit value")
	}
	return protocol.MessageCommit[RI]{
		m.Epoch,
		m.SeqNr,
		m.Signature,
	}, nil
}

func CertifiedPrepareOrCommitFromProtoMessage(m *CertifiedPrepareOrCommit) (protocol.CertifiedPrepareOrCommit, error) {
	if m == nil {
		return nil, fmt.Errorf("unable to extract a CertifiedPrepareOrCommit value")
	}
	switch poc := m.PrepareOrCommit.(type) {
	case *CertifiedPrepareOrCommit_Prepare:
		cpocp, err := certifiedPrepareFromProtoMessage(poc.Prepare)
		if err != nil {
			return nil, err
		}
		return &cpocp, nil
	case *CertifiedPrepareOrCommit_Commit:
		cpocc, err := certifiedCommitFromProtoMessage(poc.Commit)
		if err != nil {
			return nil, err
		}
		return &cpocc, nil
	default:
		return nil, fmt.Errorf("unknown case of CertifiedPrepareOrCommit")
	}
}

func certifiedPrepareFromProtoMessage(m *CertifiedPrepare) (protocol.CertifiedPrepare, error) {
	if m == nil {
		return protocol.CertifiedPrepare{}, fmt.Errorf("unable to extract a CertifiedPrepare value")
	}
	var outcomeInputsDigest protocol.OutcomeInputsDigest
	copy(outcomeInputsDigest[:], m.OutcomeInputsDigest)
	prepareQuorumCertificate := make([]protocol.AttributedPrepareSignature, 0, len(m.PrepareQuorumCertificate))
	for _, aps := range m.PrepareQuorumCertificate {
		prepareQuorumCertificate = append(prepareQuorumCertificate, protocol.AttributedPrepareSignature{
			aps.GetSignature(),
			commontypes.OracleID(aps.GetSigner()),
		})
	}
	return protocol.CertifiedPrepare{
		m.PrepareEpoch,
		m.SeqNr,
		outcomeInputsDigest,
		m.Outcome,
		prepareQuorumCertificate,
	}, nil

}

func certifiedCommitFromProtoMessage(m *CertifiedCommit) (protocol.CertifiedCommit, error) {
	if m == nil {
		return protocol.CertifiedCommit{}, fmt.Errorf("unable to extract a CertifiedCommit value")
	}
	commitQuorumCertificate := make([]protocol.AttributedCommitSignature, 0, len(m.CommitQuorumCertificate))
	for _, aps := range m.CommitQuorumCertificate {
		commitQuorumCertificate = append(commitQuorumCertificate, protocol.AttributedCommitSignature{
			aps.GetSignature(),
			commontypes.OracleID(aps.GetSigner()),
		})
	}
	return protocol.CertifiedCommit{
		m.CommitEpoch,
		m.SeqNr,
		m.Outcome,
		commitQuorumCertificate,
	}, nil
}

func signedHighestCertifiedTimestampFromProtoMessage(m *SignedHighestCertifiedTimestamp) (protocol.SignedHighestCertifiedTimestamp, error) {
	if m == nil {
		return protocol.SignedHighestCertifiedTimestamp{}, fmt.Errorf("unable to extract a SignedHighestCertifiedTimestamp value")
	}
	hct, err := highestCertifiedTimestampFromProtoMessage(m.HighestCertifiedTimestamp)
	if err != nil {
		return protocol.SignedHighestCertifiedTimestamp{}, err
	}
	return protocol.SignedHighestCertifiedTimestamp{
		hct,
		m.Signature,
	}, nil
}

func highestCertifiedTimestampFromProtoMessage(m *HighestCertifiedTimestamp) (protocol.HighestCertifiedTimestamp, error) {
	if m == nil {
		return protocol.HighestCertifiedTimestamp{}, fmt.Errorf("unable to extract a HighestCertifiedTimestamp value")
	}
	return protocol.HighestCertifiedTimestamp{
		m.SeqNr,
		m.CommittedElsePrepared,
	}, nil
}

func epochStartProofFromProtoMessage(m *EpochStartProof) (protocol.EpochStartProof, error) {
	if m == nil {
		return protocol.EpochStartProof{}, fmt.Errorf("unable to extract a EpochStartProof value")
	}
	hc, err := CertifiedPrepareOrCommitFromProtoMessage(m.HighestCertified)
	if err != nil {
		return protocol.EpochStartProof{}, err
	}
	hctqc := make([]protocol.AttributedSignedHighestCertifiedTimestamp, 0, len(m.HighestCertifiedProof))
	for _, ashct := range m.HighestCertifiedProof {
		hctqc = append(hctqc, protocol.AttributedSignedHighestCertifiedTimestamp{
			protocol.SignedHighestCertifiedTimestamp{
				protocol.HighestCertifiedTimestamp{
					ashct.GetSignedHighestCertifiedTimestamp().GetHighestCertifiedTimestamp().GetSeqNr(),
					ashct.GetSignedHighestCertifiedTimestamp().GetHighestCertifiedTimestamp().GetCommittedElsePrepared(),
				},
				ashct.GetSignedHighestCertifiedTimestamp().GetSignature(),
			},
			commontypes.OracleID(ashct.GetSigner()),
		})
	}

	return protocol.EpochStartProof{
		hc,
		hctqc,
	}, nil
}

func messageRoundStartFromProtoMessage[RI any](m *MessageRoundStart) (protocol.MessageRoundStart[RI], error) {
	if m == nil {
		return protocol.MessageRoundStart[RI]{}, fmt.Errorf("unable to extract a MessageRoundStart value")
	}
	return protocol.MessageRoundStart[RI]{
		m.Epoch,
		m.SeqNr,
		m.Query,
	}, nil
}

func messageObservationFromProtoMessage[RI any](m *MessageObservation) (protocol.MessageObservation[RI], error) {
	if m == nil {
		return protocol.MessageObservation[RI]{}, fmt.Errorf("unable to extract a MessageObservation value")
	}
	so, err := signedObservationFromProtoMessage(m.SignedObservation)
	if err != nil {
		return protocol.MessageObservation[RI]{}, err
	}
	return protocol.MessageObservation[RI]{
		m.Epoch,
		m.SeqNr,
		so,
	}, nil
}

func messageReportSignaturesFromProtoMessage[RI any](m *MessageReportSignatures) (protocol.MessageReportSignatures[RI], error) {
	if m == nil {
		return protocol.MessageReportSignatures[RI]{}, fmt.Errorf("unable to extract a MessageReportSignatures value")
	}
	return protocol.MessageReportSignatures[RI]{
		m.SeqNr,
		m.ReportSignatures,
	}, nil
}

func messageCertifiedCommitRequestFromProtoMessage[RI any](m *MessageCertifiedCommitRequest) (protocol.MessageCertifiedCommitRequest[RI], error) {
	if m == nil {
		return protocol.MessageCertifiedCommitRequest[RI]{}, fmt.Errorf("unable to extract a MessageCertifiedCommitRequest value")
	}
	return protocol.MessageCertifiedCommitRequest[RI]{
		m.SeqNr,
	}, nil
}

func messageCertifiedCommitFromProtoMessage[RI any](m *MessageCertifiedCommit) (protocol.MessageCertifiedCommit[RI], error) {
	if m == nil {
		return protocol.MessageCertifiedCommit[RI]{}, fmt.Errorf("unable to extract a MessageCertifiedCommit value")
	}
	cpocc, err := certifiedCommitFromProtoMessage(m.CertifiedCommit)
	if err != nil {
		return protocol.MessageCertifiedCommit[RI]{}, err
	}
	return protocol.MessageCertifiedCommit[RI]{
		cpocc,
	}, nil
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
		return protocol.SignedObservation{}, fmt.Errorf("unable to extract a SignedObservation value")
	}

	return protocol.SignedObservation{
		m.Observation,
		m.Signature,
	}, nil
}

func PacemakerStateFromProtoMessage(m *PacemakerState) (protocol.PacemakerState, error) {
	if m == nil {
		return protocol.PacemakerState{}, fmt.Errorf("unable to extract a PacemakerState value")
	}

	return protocol.PacemakerState{
		m.Epoch,
		m.HighestSentNewEpochWish,
	}, nil
}
