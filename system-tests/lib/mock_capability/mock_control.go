package mock_capability

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	pb2 "github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/system-tests/lib/mock_capability/pb"
)

type MockCapabilityController struct {
	lggr    zerolog.Logger
	Clients []pb.MockCapabilityClient
}

func NewMockCapabilityController(lggr zerolog.Logger) *MockCapabilityController {
	return &MockCapabilityController{Clients: make([]pb.MockCapabilityClient, 0), lggr: lggr}
}

// ConnectAll connects to all addresses, for CTFv2 test useInsecure should be true, for CRIB useInsecure should be false
func (c *MockCapabilityController) ConnectAll(addresses []string, useInsecure bool) error {
	for _, p := range addresses {
		client, err := proxyConnectToOne(p, useInsecure)
		if err != nil {
			return err
		}
		c.Clients = append(c.Clients, client)

	}
	return nil
}

func (c *MockCapabilityController) RegisterToWorkflow(ctx context.Context, info pb.RegisterToWorkflowRequest) error {
	for _, client := range c.Clients {
		_, err := client.RegisterToWorkflow(ctx, &info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MockCapabilityController) Execute(ctx context.Context, info pb.ExecutableRequest) error {
	for _, client := range c.Clients {
		_, err := client.Execute(ctx, &info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MockCapabilityController) CreateCapability(ctx context.Context, info pb.CapabilityInfo) error {
	for _, client := range c.Clients {
		_, err := client.CreateCapability(ctx, &info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MockCapabilityController) SendTrigger(ctx context.Context, id string, eventID string, payload []byte) error {
	for _, client := range c.Clients {
		data := pb.SendTriggerEventRequest{
			ID:      id,
			EventID: eventID,
			Payload: payload,
		}
		framework.L.Info().Msg(fmt.Sprintf("Sending trigger response %s:%s", id, eventID))

		_, err := client.SendTriggerEvent(ctx, &data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MockCapabilityController) HookExecutables(ctx context.Context, ch chan pb.ExecutableRequest) error {
	for _, client := range c.Clients {
		hook, errC := client.HookExecutables(context.TODO())
		if errC != nil {
			return errC
		}

		go func() {
			for {
				c.lggr.Info().Msg("Waiting for hook event")
				resp, err := hook.Recv()
				if err == io.EOF {
					c.lggr.Error().Msgf("Recieved EOF from hook %s", err)
					return
				}
				if err != nil {
					log.Fatalf("can not receive %v", err)
				}
				ch <- *resp
				c.lggr.Info().Msgf("Got hook event %v+", resp)

				//Process request
				r := pb.ExecutableResponse{
					ID:             resp.ID,
					CapabilityType: resp.CapabilityType,
					Value:          resp.Inputs,
				}
				err = hook.Send(&r)
				if err != nil {
					panic(err.Error())
				}
				c.lggr.Info().Msgf("Sent hook response %v+", r)

			}
		}()
	}
	return nil
}

func proxyConnectToOne(address string, useInsecure bool) (pb.MockCapabilityClient, error) {
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	if useInsecure {
		creds = insecure.NewCredentials()
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	client := pb.NewMockCapabilityClient(conn)
	return client, nil

}

func MapToBytes(m *values.Map) ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	pm := make(map[string]*pb2.Value)
	for k, v := range m.Underlying {
		pm[k] = values.Proto(v)
	}
	bytes, err := proto.Marshal(pb2.NewMapValue(pm))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
