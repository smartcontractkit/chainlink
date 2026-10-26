package request

// modify with care - any changes will break downstream kafka consumers.
const (
	TransmissionEventSchema   = "/workflows/v1/transmit_schedule_event.proto"
	TransmissionEventProtoPkg = "workflows.v1"
	TransmissionEventEntity   = "TransmissionsScheduledEvent"

	// ScheduleAllAtOnce is the transmission schedule used by all V2 capabilities (S = [N]):
	// every DON member transmits immediately, without staged delays. It is the value written
	// to the ScheduleType field of TransmissionsScheduledEvent.
	ScheduleAllAtOnce = "allAtOnce"
)
