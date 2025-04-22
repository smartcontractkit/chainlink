package request

// modify with care - any changes will break downstream kafka consumers.
const (
	TransmissionEventSchema = "/cre-events-workflow-started/v1"
	TransmissionEventProto  = "request"
	TransmissionEventPkg    = "TransmitScheduleEvent"
)
