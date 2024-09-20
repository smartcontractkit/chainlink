package logger

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/pb"
	"google.golang.org/protobuf/proto"
	"log"
)

// beholderCustomMessageLogger is used to emit logs as custom messages through Beholder
type beholderCustomMessageLogger struct {
	Logger
}

func NewBeholderCustomMessageLogger(l Logger) *beholderCustomMessageLogger {
	return &beholderCustomMessageLogger{l}
}

func (s *beholderCustomMessageLogger) Errorf(format string, values ...interface{}) {
	sendCustomMessage(format, values...)
	s.Logger.Errorf(format, values...)
}

func sendCustomMessage(format string, values ...interface{}) {
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
