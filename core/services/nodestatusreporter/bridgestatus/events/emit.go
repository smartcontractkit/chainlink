package events

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// EmitBridgeStatusEvent emits a Bridge Status event through the provided custmsg.MessageEmitter
func EmitBridgeStatusEvent(ctx context.Context, emitter beholder.Emitter, event *BridgeStatusEvent) error {
	if event.Timestamp == "" {
		event.Timestamp = time.Now().Format(time.RFC3339Nano)
	}

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal BridgeStatusEvent: %w", err)
	}

	err = emitter.Emit(ctx, eventBytes,
		"beholder_data_schema", SchemaBridgeStatus,
		"beholder_domain", "data-feeds",
		"beholder_entity", fmt.Sprintf("%s.%s", ProtoPkg, BridgeStatusEventEntity),
	)
	if err != nil {
		return fmt.Errorf("failed to emit BridgeStatusEvent: %w", err)
	}

	return nil
}
