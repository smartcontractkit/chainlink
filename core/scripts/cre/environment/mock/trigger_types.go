package mock

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	cron2 "github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"google.golang.org/protobuf/types/known/anypb"

	pb2 "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
)

// TriggerType defines the type of mock trigger
type TriggerType string

const (
	TriggerTypeCron TriggerType = "cron"
	// Add other trigger types as needed
)

// getTriggerRequest returns a configured SendTriggerEventRequest based on the requested type
func getTriggerRequest(triggerType TriggerType) (*pb2.SendTriggerEventRequest, error) {
	switch triggerType {
	case TriggerTypeCron:
		// First create the payload
		payload := &cron2.LegacyPayload{ //nolint:staticcheck // legacy
			ScheduledExecutionTime: time.Now().Format(time.RFC3339Nano),
		}

		// Convert it to anypb.Any
		anyPayload, err := anypb.New(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to convert payload: %w", err)
		}

		return &pb2.SendTriggerEventRequest{
			ID:      uuid.New().String(),
			Payload: anyPayload,
		}, nil
	default:
		return nil, fmt.Errorf("unknown trigger type: %s", triggerType)
	}
}
