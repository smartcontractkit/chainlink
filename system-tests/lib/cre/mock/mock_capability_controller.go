package mockcapability

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	pb2 "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
)

type Controller struct {
	lggr  zerolog.Logger
	Nodes []MockClient
}

type MockClient struct {
	API pb2.MockCapabilityClient
	URL string
}

type OCRTriggerEvent struct {
	ConfigDigest []byte
	SeqNr        uint64
	Report       []byte
	Sigs         []OCRTriggerEventSig
}

type OCRTriggerEventSig struct {
	Signature []byte
	Signer    uint32
}

func NewMockCapabilityController(lggr zerolog.Logger) *Controller {
	return &Controller{Nodes: make([]MockClient, 0), lggr: lggr}
}

func NewMockCapabilityControllerFromCache(lggr zerolog.Logger, useInsecure bool) (*Controller, error) {
	bytes, err := os.ReadFile("cache/mock-clients.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read URLs from cache: %w", err)
	}

	addresses := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	if len(addresses) == 0 {
		return nil, errors.New("no URLs found in cache file")
	}

	controller := NewMockCapabilityController(lggr)
	if err := controller.ConnectAll(addresses, useInsecure, false); err != nil {
		return nil, fmt.Errorf("failed to connect to cached URLs: %w", err)
	}

	return controller, nil
}

// ConnectAll connects to all addresses, for CTFv2 test useInsecure should be true, for CRIB useInsecure should be false
func (c *Controller) ConnectAll(addresses []string, useInsecure bool, cacheClients bool) error {
	if cacheClients {
		cacheDir := "cache"
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return fmt.Errorf("failed to create cache directory: %w", err)
		}

		urlsBytes := []byte(strings.Join(addresses, "\n"))
		if err := os.WriteFile("cache/mock-clients.txt", urlsBytes, 0600); err != nil {
			return fmt.Errorf("failed to save URLs to cache: %w", err)
		}
	}
	for _, p := range addresses {
		client, err := proxyConnectToOne(p, useInsecure)
		if err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, client)
	}

	return nil
}

func (c *Controller) RegisterToWorkflow(ctx context.Context, info *pb2.RegisterToWorkflowRequest) error {
	for _, client := range c.Nodes {
		_, err := client.API.RegisterToWorkflow(ctx, info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) Execute(ctx context.Context, info *pb2.ExecutableRequest) error {
	for _, client := range c.Nodes {
		_, err := client.API.Execute(ctx, info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) CreateCapability(ctx context.Context, info *pb2.CapabilityInfo) error {
	for _, client := range c.Nodes {
		_, err := client.API.CreateCapability(ctx, info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) SendTrigger(ctx context.Context, id string, eventID string, payload []byte) error {
	for _, client := range c.Nodes {
		data := pb2.SendTriggerEventRequest{
			ID:      id,
			EventID: eventID,
			Payload: payload,
		}

		framework.L.Info().Msg(fmt.Sprintf("Sending trigger response %s:%s", id, eventID))

		_, err := client.API.SendTriggerEvent(ctx, &data)
		if err != nil {
			return err
		}
	}
	return nil
}

type CapInfos struct {
	Node         string
	Capabilities []capabilities.CapabilityInfo
}

func (c *Controller) List(ctx context.Context) ([]CapInfos, error) {
	info := make([]CapInfos, 0)
	for _, client := range c.Nodes {
		data, err := client.API.List(ctx, &pb2.ListRequest{})
		if err != nil {
			return nil, err
		}
		framework.L.Info().Msgf("Fetching capabilityes for node %s", client.URL)
		caps := make([]capabilities.CapabilityInfo, 0)
		for _, d := range data.CapInfos {
			caps = append(caps, capabilities.CapabilityInfo{
				ID:             d.ID,
				CapabilityType: capabilities.CapabilityType(d.CapabilityType),
				Description:    d.Description,
				IsLocal:        d.IsLocal,
			})
		}

		info = append(info, CapInfos{
			Node:         client.URL,
			Capabilities: caps,
		})
	}
	return info, nil
}

func (c *Controller) HookExecutables(ctx context.Context, ch chan capabilities.CapabilityRequest) error {
	for _, client := range c.Nodes {
		hook, errC := client.API.HookExecutables(context.TODO())
		if errC != nil {
			return fmt.Errorf("cannot hook into executable at %s: %w", client.URL, errC)
		}

		go func() {
			for {
				c.lggr.Info().Msg("Waiting for hook event")
				resp, err := hook.Recv()
				if errors.Is(err, io.EOF) {
					c.lggr.Error().Msgf("Received EOF from hook %s", err)
					return
				}
				if err != nil {
					log.Fatalf("can not receive %v", err)
				}

				config, err := BytesToMap(resp.Config)
				if err != nil {
					log.Fatalf("can not decode config: %v", err)
				}
				input, err := BytesToMap(resp.Inputs)
				if err != nil {
					log.Fatalf("can not decode input: %v", err)
				}

				c.lggr.Info().Msgf("Got hook event %s", resp.ID)
				ch <- capabilities.CapabilityRequest{
					Metadata: capabilities.RequestMetadata{
						WorkflowID:               resp.RequestMetadata.WorkflowID,
						WorkflowOwner:            resp.RequestMetadata.WorkflowOwner,
						WorkflowExecutionID:      resp.RequestMetadata.WorkflowExecutionID,
						WorkflowName:             resp.RequestMetadata.WorkflowName,
						WorkflowDonID:            resp.RequestMetadata.WorkflowDonID,
						WorkflowDonConfigVersion: resp.RequestMetadata.WorkflowDonConfigVersion,
						ReferenceID:              resp.RequestMetadata.ReferenceID,
						DecodedWorkflowName:      resp.RequestMetadata.DecodedWorkflowName,
					},
					Config: config,
					Inputs: input,
				}
				c.lggr.Info().Msgf("Got hook event %s", resp.ID)

				r := pb2.ExecutableResponse{
					ID:             resp.ID,
					CapabilityType: resp.CapabilityType,
					Value:          resp.Inputs,
				}
				err = hook.Send(&r)
				if err != nil {
					panic(err.Error())
				}
			}
		}()
	}
	return nil
}

func proxyConnectToOne(address string, useInsecure bool) (MockClient, error) {
	//nolint:gosec // disable G402
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	if useInsecure {
		creds = insecure.NewCredentials()
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return MockClient{}, err
	}
	client := pb2.NewMockCapabilityClient(conn)
	return MockClient{API: client, URL: address}, nil
}
