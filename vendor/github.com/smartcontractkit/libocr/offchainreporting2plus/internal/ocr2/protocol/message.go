package protocol

import (
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// EventToPacemaker is the interface used to pass in-process events to the
// leader-election protocol.
type EventToPacemaker interface {
	// processPacemaker is called when the local oracle process invokes an event
	// intended for the leader-election protocol.
	processPacemaker(pace *pacemakerState)
}

// EventProgress is used to process the "progress" event passed by the local
// oracle from its the reporting protocol to the leader-election protocol. It is
// sent by the reporting protocol when the leader has produced a valid new
// report.
type EventProgress struct{}

var _ EventToPacemaker = (*EventProgress)(nil) // implements EventToPacemaker

func (ev EventProgress) processPacemaker(pace *pacemakerState) {
	pace.eventProgress()
}

// EventChangeLeader is used to process the "change-leader" event passed by the
// local oracle from its the reporting protocol to the leader-election protocol
type EventChangeLeader struct{}

var _ EventToPacemaker = (*EventChangeLeader)(nil) // implements EventToPacemaker

func (ev EventChangeLeader) processPacemaker(pace *pacemakerState) {
	pace.eventChangeLeader()
}

type EventToReportFinalization interface {
	processReportFinalization(repfin *reportFinalizationState)
}

// EventToTransmission is the interface used to pass a completed report to the
// protocol which will transmit it to the on-chain smart contract.
type EventToTransmission interface {
	processTransmission(t *transmissionState)
}

// Message is the interface used to pass an inter-oracle message to the local
// oracle process.
//

type Message interface {
	// CheckSize checks whether the given message conforms to the limits imposed by
	// reportingPluginLimits
	CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool

	// process passes this Message instance to the oracle o, as a message from
	// oracle with the given sender index
	process(o *oracleState, sender commontypes.OracleID)
}

// MessageWithSender records a msg with the index of the sender oracle
type MessageWithSender struct {
	Msg    Message
	Sender commontypes.OracleID
}

// MessageToPacemaker is the interface used to pass a message to the local
// leader-election protocol
type MessageToPacemaker interface {
	Message

	// process passes this MessageToPacemaker instance to the oracle o, as a
	// message from oracle with the given sender index
	processPacemaker(pace *pacemakerState, sender commontypes.OracleID)
}

// MessageToPacemakerWithSender records a msg with the idx of the sender oracle
type MessageToPacemakerWithSender struct {
	msg    MessageToPacemaker
	sender commontypes.OracleID
}

// MessageToReportGeneration is the interface used to pass an inter-oracle message
// to the local oracle reporting process.
type MessageToReportGeneration interface {
	Message

	// processReportGeneration is called to send this message to the local oracle
	// reporting process.
	processReportGeneration(repgen *reportGenerationState, sender commontypes.OracleID)

	epoch() uint32
}

// MessageToReportGenerationWithSender records a message destined for the oracle
// reporting
type MessageToReportGenerationWithSender struct {
	msg    MessageToReportGeneration
	sender commontypes.OracleID
}

type MessageToReportFinalization interface {
	Message

	processReportFinalization(repfin *reportFinalizationState, sender commontypes.OracleID)
}

type MessageToReportFinalizationWithSender struct {
	msg    MessageToReportFinalization
	sender commontypes.OracleID
}

// MessageNewEpoch corresponds to the "newepoch(epoch_number)" message from alg.
// 1. It indicates that the node believes the protocol should move to the
// specified epoch.
type MessageNewEpoch struct {
	Epoch uint32
}

var _ MessageToPacemaker = (*MessageNewEpoch)(nil)

func (msg MessageNewEpoch) CheckSize(types.ReportingPluginLimits) bool {
	return true
}

func (msg MessageNewEpoch) process(o *oracleState, sender commontypes.OracleID) {
	o.chNetToPacemaker <- MessageToPacemakerWithSender{msg, sender}
}

func (msg MessageNewEpoch) processPacemaker(pace *pacemakerState, sender commontypes.OracleID) {
	pace.messageNewepoch(msg, sender)
}

// MessageObserveReq corresponds to the "observe-req" message from alg. 2. The
// leader transmits this to request observations from participating oracles, so
// that it can collate them into a report.
type MessageObserveReq struct {
	Epoch uint32
	Round uint8
	Query types.Query
}

var _ MessageToReportGeneration = (*MessageObserveReq)(nil)

func (msg MessageObserveReq) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return len(msg.Query) <= reportingPluginLimits.MaxQueryLength
}

func (msg MessageObserveReq) process(o *oracleState, sender commontypes.OracleID) {
	o.reportGenerationMessage(msg, sender)
}

func (msg MessageObserveReq) processReportGeneration(repgen *reportGenerationState, sender commontypes.OracleID) {
	repgen.messageObserveReq(msg, sender)
}

func (msg MessageObserveReq) epoch() uint32 {
	return msg.Epoch
}

// MessageObserve corresponds to the "observe" message from alg. 2.
// Participating oracles send this back to the leader in response to
// MessageObserveReq's.
type MessageObserve struct {
	Epoch             uint32
	Round             uint8
	SignedObservation SignedObservation
}

var _ MessageToReportGeneration = (*MessageObserve)(nil)

func (msg MessageObserve) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return len(msg.SignedObservation.Observation) <= reportingPluginLimits.MaxObservationLength
}

func (msg MessageObserve) process(o *oracleState, sender commontypes.OracleID) {
	o.reportGenerationMessage(msg, sender)
}

func (msg MessageObserve) processReportGeneration(repgen *reportGenerationState, sender commontypes.OracleID) {
	repgen.messageObserve(msg, sender)
}

func (msg MessageObserve) epoch() uint32 {
	return msg.Epoch
}

// MessageReportReq corresponds to the "report-req" message from alg. 2. It is
// sent by the epoch leader with collated observations for the participating
// oracles to sign.
type MessageReportReq struct {
	Epoch                        uint32
	Round                        uint8
	Query                        types.Query
	AttributedSignedObservations []AttributedSignedObservation
}

func (msg MessageReportReq) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	if !(len(msg.Query) <= reportingPluginLimits.MaxQueryLength) {
		return false
	}

	if !(len(msg.AttributedSignedObservations) <= types.MaxOracles) {
		return false
	}

	for _, aso := range msg.AttributedSignedObservations {
		if !(len(aso.SignedObservation.Observation) <= reportingPluginLimits.MaxObservationLength) {
			return false
		}
	}

	return true
}

func (msg MessageReportReq) process(o *oracleState, sender commontypes.OracleID) {
	o.reportGenerationMessage(msg, sender)
}

func (msg MessageReportReq) processReportGeneration(repgen *reportGenerationState, sender commontypes.OracleID) {
	repgen.messageReportReq(msg, sender)
}

func (msg MessageReportReq) epoch() uint32 {
	return msg.Epoch
}

var _ MessageToReportGeneration = (*MessageReportReq)(nil)

// MessageReport corresponds to the "report" message from alg. 2. It is sent by
// participating oracles in response to a MessageReportReq, and contains the
// final form of the report, based on the collated observations, and the sending
// oracle's signature.
type MessageReport struct {
	Epoch          uint32
	Round          uint8
	AttestedReport AttestedReportOne
}

var _ MessageToReportGeneration = (*MessageReport)(nil)

func (msg MessageReport) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return len(msg.AttestedReport.Report) <= reportingPluginLimits.MaxReportLength
}

func (msg MessageReport) process(o *oracleState, sender commontypes.OracleID) {
	o.reportGenerationMessage(msg, sender)
}

func (msg MessageReport) processReportGeneration(repgen *reportGenerationState, sender commontypes.OracleID) {
	repgen.messageReport(msg, sender)
}

func (msg MessageReport) epoch() uint32 {
	return msg.Epoch
}

// MessageFinal corresponds to the "final" message in alg. 2. It is sent by the
// current leader with the aggregated signature(s) to all participating oracles,
// for them to participate in the subsequent transmission of the report to the
// on-chain contract.
type MessageFinal struct {
	Epoch          uint32
	Round          uint8
	H              [32]byte
	AttestedReport AttestedReportMany
}

var _ MessageToReportGeneration = (*MessageFinal)(nil)

func (msg MessageFinal) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return len(msg.AttestedReport.AttributedSignatures) <= types.MaxOracles &&
		len(msg.AttestedReport.Report) <= reportingPluginLimits.MaxReportLength
}

func (msg MessageFinal) process(o *oracleState, sender commontypes.OracleID) {
	o.reportGenerationMessage(msg, sender)
}

func (msg MessageFinal) processReportGeneration(repgen *reportGenerationState, sender commontypes.OracleID) {
	repgen.messageFinal(msg, sender)
}

func (msg MessageFinal) epoch() uint32 {
	return msg.Epoch
}

var _ MessageToReportFinalization = MessageFinalEcho{}

// MessageFinalEcho corresponds to the "final-echo" message in alg. 2. It is
// broadcast by all oracles to all other oracles, to ensure that all can play
// their role in transmitting the report to the on-chain contract.
type MessageFinalEcho struct {
	MessageFinal
}

func (msg MessageFinalEcho) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return msg.MessageFinal.CheckSize(reportingPluginLimits)
}

func (msg MessageFinalEcho) process(o *oracleState, sender commontypes.OracleID) {
	o.chNetToReportFinalization <- MessageToReportFinalizationWithSender{msg, sender}
}

func (msg MessageFinalEcho) processReportFinalization(repfin *reportFinalizationState, sender commontypes.OracleID) {
	repfin.messageFinalEcho(msg, sender)
}

type EventFinal struct {
	MessageFinal
}

var _ EventToReportFinalization = EventFinal{} // implements EventToReportFinalization

func (ev EventFinal) processReportFinalization(repfin *reportFinalizationState) {
	repfin.eventFinal(ev)
}

// EventTransmit is used to process the "transmit" event passed by the local
// reporting protocol to to the local transmit-to-the-onchain-smart-contract
// protocol.
type EventTransmit struct {
	Epoch          uint32
	Round          uint8
	H              [32]byte
	AttestedReport AttestedReportMany
}

var _ EventToTransmission = (*EventTransmit)(nil) // implements EventToTransmission

func (ev EventTransmit) processTransmission(t *transmissionState) {
	t.eventTransmit(ev)
}
