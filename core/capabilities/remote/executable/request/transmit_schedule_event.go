package request

// modify with care - any changes will break downstream kafka consumers.
const (
	TransmissionEventSchema   = "github.com/smartcontractkit/chainlink-protos/workflows/workflows/v1/transmit_schedule_event.proto"
	TransmissionEventProtoPkg = "workflows.v1"
	TransmissionEventEntity   = "TransmitScheduleEvent"
)
