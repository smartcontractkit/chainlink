package serialization

import (
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/protocol"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/signature"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func Serialize(m protocol.Message) (b []byte, err error) {
	pmsg, err := toProtoMessage(m)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(pmsg)
}

func Deserialize(b []byte) (protocol.Message, error) {
	msgWrapper := &MessageWrapper{}
	err := proto.Unmarshal(b, msgWrapper)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal protobuf")
	}
	return msgWrapper.fromProtoMessage()
}

func toProtoMessage(m protocol.Message) (proto.Message, error) {
	msgWrapper := MessageWrapper{}
	switch v := m.(type) {
	case protocol.MessageNewEpoch:
		pm := &MessageNewEpoch{
			Epoch: uint64(v.Epoch),
		}
		msgWrapper.Msg = &MessageWrapper_MessageNewEpoch{pm}
	case protocol.MessageObserveReq:
		pm := &MessageObserveReq{
			Round: uint64(v.Round),
			Epoch: uint64(v.Epoch),
		}
		msgWrapper.Msg = &MessageWrapper_MessageObserveReq{pm}
	case protocol.MessageObserve:
		pm := &MessageObserve{
			Round: uint64(v.Round),
			Epoch: uint64(v.Epoch),
			Obs:   observationToProtoMessage(v.Obs),
		}
		msgWrapper.Msg = &MessageWrapper_MessageObserve{pm}
	case protocol.MessageReportReq:
		pm := &MessageReportReq{
			Round: uint64(v.Round),
			Epoch: uint64(v.Epoch),
		}
		for _, o := range v.Observations {
			pm.Observations = append(pm.Observations, observationToProtoMessage(o))
		}
		msgWrapper.Msg = &MessageWrapper_MessageReportReq{pm}
	case protocol.MessageReport:
		pm := &MessageReport{
			Epoch:          uint64(v.Epoch),
			Round:          uint64(v.Round),
			ContractReport: contractReportToProtoMessage(v.ContractReport),
		}
		msgWrapper.Msg = &MessageWrapper_MessageReport{pm}
	case protocol.MessageFinal:
		msgWrapper.Msg = &MessageWrapper_MessageFinal{finalToProtoMessage(v)}
	case protocol.MessageFinalEcho:
		msgWrapper.Msg = &MessageWrapper_MessageFinalEcho{
			&MessageFinalEcho{Final: finalToProtoMessage(v.MessageFinal)},
		}
	default:
		return nil, errors.Errorf("Unable to serialize message of type %T", m)

	}
	return &msgWrapper, nil
}

func typesObservationToProtoMessage(o observation.Observation) *ObservationValue {
	return &ObservationValue{Value: o.Marshal()}
}

func observationToProtoMessage(o protocol.Observation) *Observation {
	sig := o.Sig
	if sig == nil {
		sig = []byte{}
	}
	return &Observation{
		Ctx:       reportingContextToProtoMessage(o.Ctx),
		OracleID:  uint32(o.OracleID),
		Value:     typesObservationToProtoMessage(o.Value),
		Signature: sig,
	}
}

func reportingContextToProtoMessage(r signature.ReportingContext) *ReportingContext {
	return &ReportingContext{
		ConfigDigest: r.ConfigDigest[:],
		Epoch:        uint64(r.Epoch),
		Round:        uint64(r.Round),
	}
}

func (wrapper *MessageWrapper) fromProtoMessage() (protocol.Message, error) {
	switch msg := wrapper.Msg.(type) {
	case *MessageWrapper_MessageNewEpoch:
		return wrapper.GetMessageNewEpoch().fromProtoMessage()
	case *MessageWrapper_MessageObserveReq:
		return wrapper.GetMessageObserveReq().fromProtoMessage()
	case *MessageWrapper_MessageObserve:
		return wrapper.GetMessageObserve().fromProtoMessage()
	case *MessageWrapper_MessageReportReq:
		return wrapper.GetMessageReportReq().fromProtoMessage()
	case *MessageWrapper_MessageReport:
		return wrapper.GetMessageReport().fromProtoMessage()
	case *MessageWrapper_MessageFinal:
		return wrapper.GetMessageFinal().fromProtoMessage()
	case *MessageWrapper_MessageFinalEcho:
		return wrapper.GetMessageFinalEcho().fromProtoMessage()
	default:
		return nil, errors.Errorf("Unrecognised Msg type %T", msg)
	}
}

func (m *MessageNewEpoch) fromProtoMessage() (protocol.MessageNewEpoch, error) {
	return protocol.MessageNewEpoch{
		Epoch: uint32(m.Epoch),
	}, nil
}

func (m *MessageObserveReq) fromProtoMessage() (protocol.MessageObserveReq, error) {
	return protocol.MessageObserveReq{
		Round: uint8(m.Round),
		Epoch: uint32(m.Epoch),
	}, nil
}

func (m *MessageObserve) fromProtoMessage() (protocol.MessageObserve, error) {
	obs, err := m.Obs.fromProtoMessage()
	if err != nil {
		return protocol.MessageObserve{}, nil
	}
	return protocol.MessageObserve{
		Epoch: uint32(m.Epoch),
		Round: uint8(m.Round),
		Obs:   obs,
	}, nil
}

func (m *MessageReportReq) fromProtoMessage() (protocol.MessageReportReq, error) {
	observations, err := Observations(m.Observations).fromProtoMessage()
	if err != nil {
		return protocol.MessageReportReq{}, err
	}
	return protocol.MessageReportReq{
		Round:        uint8(m.Round),
		Epoch:        uint32(m.Epoch),
		Observations: observations,
	}, nil
}

func (o *ObservationValue) fromProtoMessage() (observation.Observation, error) {
	obs, err := observation.UnmarshalObservation(o.Value)
	if err != nil {
		return observation.Observation{}, errors.Errorf(`could not deserialize bytes as `+
			`observation.Observation: "%v" from 0x%x`, err, o.Value)
	}
	return obs, nil
}

func (m *ContractReport) fromProtoMessage() (protocol.ContractReport, error) {
	if m == nil {
		return protocol.ContractReport{}, nil
	}
	values := make([]protocol.OracleValue, len(m.Values))
	for i, v := range m.Values {
		val, err := v.Value.fromProtoMessage()
		if err != nil {
			return protocol.ContractReport{}, err
		}
		values[i] = protocol.OracleValue{
			ID:    types.OracleID(v.OracleID),
			Value: val,
		}
	}
	sig := m.Sig
	if sig == nil {
		sig = []byte{}
	}
	ctx, err := m.Ctx.fromProtoMessage()
	if err != nil {
		return protocol.ContractReport{}, err
	}
	return protocol.ContractReport{
		Ctx:    ctx,
		Values: values,
		Sig:    sig,
	}, nil
}

func (r *ReportingContext) fromProtoMessage() (signature.ReportingContext, error) {
	return signature.ReportingContext{
		ConfigDigest: types.BytesToConfigDigest(r.ConfigDigest),
		Epoch:        uint32(r.Epoch),
		Round:        uint8(r.Round),
	}, nil
}

func (m *MessageReport) fromProtoMessage() (protocol.MessageReport, error) {
	contractReport, err := m.ContractReport.fromProtoMessage()
	if err != nil {
		return protocol.MessageReport{}, err
	}

	return protocol.MessageReport{
		Epoch:          uint32(m.Epoch),
		Round:          uint8(m.Round),
		ContractReport: contractReport,
	}, nil
}

func (m *ContractReportWithSignatures) fromProtoMessage() (protocol.ContractReportWithSignatures, error) {
	signatures := make([][]byte, len(m.Signatures))
	for i, s := range m.Signatures {
		sig := s.Signature
		if sig == nil {
			sig = []byte{}
		}
		signatures[i] = sig
	}
	contractReport, err := m.ContractReport.fromProtoMessage()
	if err != nil {
		return protocol.ContractReportWithSignatures{}, err
	}
	return protocol.ContractReportWithSignatures{
		ContractReport: contractReport,
		Signatures:     signatures,
	}, nil
}

func (m *MessageFinal) fromProtoMessage() (protocol.MessageFinal, error) {
	report, err := m.Report.fromProtoMessage()
	if err != nil {
		return protocol.MessageFinal{}, nil
	}
	return protocol.MessageFinal{
		Epoch:  uint32(m.Epoch),
		Round:  uint8(m.Round),
		Report: report,
	}, nil
}

func (m *MessageFinalEcho) fromProtoMessage() (protocol.MessageFinalEcho, error) {
	final, err := m.Final.fromProtoMessage()
	if err != nil {
		return protocol.MessageFinalEcho{}, err
	}
	return protocol.MessageFinalEcho{MessageFinal: final}, nil
}

type Observations []*Observation

func (ms Observations) fromProtoMessage() ([]protocol.Observation, error) {
	observations := make([]protocol.Observation, len(ms))
	for i, o := range ms {
		obs, err := o.fromProtoMessage()
		if err != nil {
			return nil, err
		}
		observations[i] = obs
	}
	return observations, nil
}

func (m *Observation) fromProtoMessage() (protocol.Observation, error) {
	sig := m.Signature
	if sig == nil {
		sig = []byte{}
	}
	v, err := m.Value.fromProtoMessage()
	if err != nil {
		return protocol.Observation{}, err
	}
	ctx, err := m.Ctx.fromProtoMessage()
	if err != nil {
		return protocol.Observation{}, err
	}
	return protocol.Observation{
		Ctx:      ctx,
		Value:    v,
		Sig:      sig,
		OracleID: types.OracleID(m.OracleID),
	}, nil
}

func contractReportToProtoMessage(v protocol.ContractReport) *ContractReport {
	sig := v.Sig
	if sig == nil {
		sig = []byte{}
	}
	pm := &ContractReport{
		Ctx:    reportingContextToProtoMessage(v.Ctx),
		Sig:    sig,
		Values: make([]*OracleValue, len(v.Values)),
	}
	for i, val := range v.Values {
		pm.Values[i] = &OracleValue{
			OracleID: uint32(val.ID),
			Value:    &ObservationValue{Value: val.Value.Marshal()},
		}
	}
	return pm
}

func finalToProtoMessage(v protocol.MessageFinal) *MessageFinal {
	pm := &MessageFinal{
		Epoch: uint64(v.Epoch),
		Round: uint64(v.Round),
		Report: &ContractReportWithSignatures{
			ContractReport: contractReportToProtoMessage(v.Report.ContractReport),
			Signatures:     make([]*Signature, len(v.Report.Signatures)),
		},
	}
	for i, sig := range v.Report.Signatures {
		if sig == nil {
			sig = []byte{}
		}
		pm.Report.Signatures[i] = &Signature{Signature: sig}
	}
	return pm
}
