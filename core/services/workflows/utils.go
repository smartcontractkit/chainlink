package workflows

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"log"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/pb"
)

const WORKFLOW_ID = "WorkflowID"
const WORKFLOW_EXECUTION_ID = "WorkflowExecutionID"

func sendLogAsCustomMessage(ctx context.Context, format string, values ...interface{}) {
	keystoneLabels := []string{
		WORKFLOW_ID,
		WORKFLOW_EXECUTION_ID,
		// etc...
	}

	msg := fmt.Sprintf(format, values...)

	for _, label := range keystoneLabels {
		val := ctx.Value(label)
		if val != nil {
			msg = fmt.Sprintf("%v.%v", label, msg)
		}
	}

	// Define a custom protobuf payload to emit
	// TODO: add a generalized custom message while beholder can't emit logs
	payload := &pb.TestCustomMessage{
		StringVal: msg,
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal protobuf")
	}

	err = beholder.GetEmitter().Emit(context.Background(), payloadBytes,
		"beholder_data_schema", "/custom-message/versions/1", // required
		"beholder_data_type", "custom_message",
	)
	if err != nil {
		log.Printf("Error emitting message: %v", err)
	}
}

func sendLogAndCustomMessage(lggr logger.SugaredLogger, format string, values ...interface{}) {
	lggr.Errorf(format, values...)
	// Define a custom protobuf payload to emit
	// TODO: add a generalized custom message while beholder can't emit logs

	// logger option A - extract already applied labels from a logger

	// logger option B - wrap a logger implementation to "print" to a passed in ptr string var
	payload := &pb.TestCustomMessage{
		StringVal: fmt.Sprintf(format, values...),
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal protobuf")
	}

	err = beholder.GetEmitter().Emit(context.Background(), payloadBytes,
		"beholder_data_schema", "/custom-message/versions/1", // required
		"beholder_data_type", "custom_message",
	)
	if err != nil {
		log.Printf("Error emitting message: %v", err)
	}
}
