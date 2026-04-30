package events

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

func EmitJobSpecEvent(ctx context.Context, emitter beholder.Emitter, event *JobSpecEvent) error {
	if event.Timestamp == "" {
		event.Timestamp = time.Now().Format(time.RFC3339Nano)
	}

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal JobSpecEvent: %w", err)
	}

	err = emitter.Emit(ctx, eventBytes,
		"partitionkey", event.ExternalJobId,
		"beholder_data_schema", SchemaJobSpec,
		"beholder_domain", BeholderDomain,
		"beholder_entity", fmt.Sprintf("%s.%s", ProtoPkg, JobSpecEventEntity),
	)
	if err != nil {
		return fmt.Errorf("failed to emit JobSpecEvent: %w", err)
	}

	return nil
}
