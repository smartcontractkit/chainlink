package protocol

import "bytes"

// The functions in this file are only used in tests, hence
// the name "TestEqual" to make that more clear

func (msg MessageObserveReq) TestEqual(msg2 MessageObserveReq) bool {
	return msg.Epoch == msg2.Epoch &&
		msg.Round == msg2.Round &&
		bytes.Equal(msg.Query, msg2.Query)
}

func (msg MessageObserve) TestEqual(msg2 MessageObserve) bool {
	return msg.Epoch == msg2.Epoch &&
		msg.Round == msg2.Round &&
		msg.SignedObservation.Equal(msg2.SignedObservation)
}

func (msg MessageReportReq) TestEqual(msg2 MessageReportReq) bool {
	if !(msg.Epoch == msg2.Epoch &&
		msg.Round == msg2.Round &&
		bytes.Equal(msg.Query, msg2.Query)) {
		return false
	}
	if len(msg.AttributedSignedObservations) != len(msg2.AttributedSignedObservations) {
		return false
	}
	for i := range msg.AttributedSignedObservations {
		if !msg.AttributedSignedObservations[i].Equal(msg2.AttributedSignedObservations[i]) {
			return false
		}
	}
	return true
}

func (msg MessageReport) TestEqual(m2 MessageReport) bool {
	return msg.Epoch == m2.Epoch && msg.Round == m2.Round && msg.AttestedReport.TestEqual(m2.AttestedReport)
}

func (msg MessageFinal) TestEqual(m2 MessageFinal) bool {
	return msg.Epoch == m2.Epoch && msg.Round == m2.Round && msg.AttestedReport.TestEqual(m2.AttestedReport)
}

func (msg MessageFinalEcho) TestEqual(m2 MessageFinalEcho) bool {
	return msg.MessageFinal.TestEqual(m2.MessageFinal)
}
