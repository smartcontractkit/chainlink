package protocol

import (
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type EventToPacemaker interface {
	processPacemaker(pace *pacemakerState)
}

type EventProgress struct{}

var _ EventToPacemaker = (*EventProgress)(nil)

func (ev EventProgress) processPacemaker(pace *pacemakerState) {
	pace.eventProgress()
}

type EventChangeLeader struct{}

var _ EventToPacemaker = (*EventChangeLeader)(nil)

func (ev EventChangeLeader) processPacemaker(pace *pacemakerState) {
	pace.eventChangeLeader()
}

type EventToTransmission interface {
	processTransmission(t *transmissionState)
}

type Message interface {
	process(o *oracleState, sender types.OracleID)
}

type MessageWithSender struct {
	Msg    Message
	Sender types.OracleID
}

type MessageWithDestination struct {
	msg  Message
	dest types.OracleID
}

type MessageToPacemaker interface {
	Message

	processPacemaker(pace *pacemakerState, sender types.OracleID)
}

type MessageToPacemakerWithSender struct {
	msg    MessageToPacemaker
	sender types.OracleID
}

type MessageToReportGeneration interface {
	Message

	processReportGeneration(repgen *reportGenerationState, sender types.OracleID)
}

type MessageToReportGenerationWithSender struct {
	msg    MessageToReportGeneration
	sender types.OracleID
}

type MessageNewEpoch struct {
	Epoch uint32
}

var _ MessageToPacemaker = (*MessageNewEpoch)(nil)

func (msg MessageNewEpoch) process(o *oracleState, sender types.OracleID) {
	o.chNetToPacemaker <- MessageToPacemakerWithSender{msg, sender}
}

func (msg MessageNewEpoch) processPacemaker(pace *pacemakerState, sender types.OracleID) {
	pace.messageNewepoch(msg, sender)
}

type MessageObserveReq struct {
	Epoch uint32
	Round uint8
}

var _ MessageToReportGeneration = (*MessageObserveReq)(nil)

func (msg MessageObserveReq) process(o *oracleState, sender types.OracleID) {
	if o.chNetToReportGeneration == nil {
		panic("nil channel! o.chNetToReportGeneration")
	}
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender{msg, sender}
}

func (msg MessageObserveReq) processReportGeneration(repgen *reportGenerationState, sender types.OracleID) {
	repgen.messageObserveReq(msg, sender)
}

type MessageObserve struct {
	Epoch uint32
	Round uint8
	Obs   Observation
}

var _ MessageToReportGeneration = (*MessageObserve)(nil)

func (msg MessageObserve) process(o *oracleState, sender types.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender{msg, sender}
}

func (msg MessageObserve) processReportGeneration(repgen *reportGenerationState, sender types.OracleID) {
	repgen.messageObserve(msg, sender)
}

func (msg MessageObserve) Equal(msg2 MessageObserve) bool {
	return msg.Epoch == msg2.Epoch &&
		msg.Round == msg2.Round &&
		msg.Obs.Equal(msg2.Obs)
}

type MessageReportReq struct {
	Epoch        uint32
	Round        uint8
	Observations []Observation
}

func (msg MessageReportReq) process(o *oracleState, sender types.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender{msg, sender}
}

func (msg MessageReportReq) processReportGeneration(repgen *reportGenerationState, sender types.OracleID) {
	repgen.messageReportReq(msg, sender)
}

var _ MessageToReportGeneration = (*MessageReportReq)(nil)

type MessageReport struct {
	Epoch          uint32
	Round          uint8
	ContractReport ContractReport
}

var _ MessageToReportGeneration = (*MessageReport)(nil)

func (msg MessageReport) process(o *oracleState, sender types.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender{msg, sender}
}

func (msg MessageReport) processReportGeneration(repgen *reportGenerationState, sender types.OracleID) {
	repgen.messageReport(msg, sender)
}

func (msg MessageReport) Equals(m2 MessageReport) bool {
	return msg.Epoch == m2.Epoch && msg.Round == m2.Round && msg.ContractReport.Equal(m2.ContractReport)
}

type MessageFinal struct {
	Epoch  uint32
	Leader types.OracleID
	Round  uint8
	Report ContractReportWithSignatures
}

var _ MessageToReportGeneration = (*MessageFinal)(nil)

func (msg MessageFinal) process(o *oracleState, sender types.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender{msg, sender}
}

func (msg MessageFinal) processReportGeneration(repgen *reportGenerationState, sender types.OracleID) {
	repgen.messageFinal(msg, sender)
}

func (msg MessageFinal) Equals(m2 MessageFinal) bool {
	return msg.Epoch == m2.Epoch && msg.Round == m2.Round && msg.Report.Equals(m2.Report)
}

type MessageFinalEcho struct {
	MessageFinal
}

func (msg MessageFinalEcho) process(o *oracleState, sender types.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender{msg, sender}
}

func (msg MessageFinalEcho) processReportGeneration(repgen *reportGenerationState, sender types.OracleID) {
	repgen.messageFinalEcho(msg, sender)
}

func (msg MessageFinalEcho) Equals(m2 MessageFinalEcho) bool {
	return msg.MessageFinal.Equals(m2.MessageFinal)
}

type EventTransmit struct{ ContractReportWithSignatures }

var _ EventToTransmission = (*EventTransmit)(nil)

func (ev EventTransmit) processTransmission(t *transmissionState) {
	t.eventTransmit(ev)
}
