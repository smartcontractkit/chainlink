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

func logAndSendCustomMessage(lggr logger.Logger, format string, values ...interface{}) {
	lggr.Errorf(format, values...)
	// Define a custom protobuf payload to emit
	// TODO: add a generalized custom message while beholder can't emit logs
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
