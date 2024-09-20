package beholder

import (
	"context"
	"log"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/beholder/pb"
)

func init() {

}

func Emit(participant, role, model, component, event string) error {
	beholderConfig := beholder.TestDefaultConfig()
	// Bootstrap Beholder Client
	client, err := beholder.NewClient(beholderConfig)
	if err != nil {
		log.Println(err.Error())
	}

	beholder.SetClient(client)

	tm := time.Now().Format("2006-01-02 15:04:05")

	// Define a custom protobuf payload to emit
	payload := &pb.Event{
		Participant: participant,
		Role:        role,
		Model:       model,
		Component:   component,
		Event:       event,
		Timestamp:   tm,
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return err
	}

	return beholder.GetEmitter().Emit(context.Background(), payloadBytes,
		"beholder_data_schema", "/event/versions/1", // required
		"beholder_data_type", "custom_message",
		"package_name", "capabilities_test",
		"participant", participant,
		"role", role,
		"model", model,
		"component", component,
		"event", event,
		"timestamp", tm,
	)
}
