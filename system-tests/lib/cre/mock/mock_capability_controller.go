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
	"time"

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

func (c *Controller) SendTrigger(ctx context.Context, message *pb2.SendTriggerEventRequest) error {
	for _, client := range c.Nodes {
		framework.L.Info().Msg(fmt.Sprintf("Sending trigger event %s to subscribers of %s", message.ID, message.TriggerID))

		_, err := client.API.SendTriggerEvent(ctx, message)
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
				c.lggr.Info().Msg("Waiting for execute events")
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

				c.lggr.Info().Msgf("Got execute event for %s with workflowID %s, executionID %s", resp.ID, resp.RequestMetadata.WorkflowID, resp.RequestMetadata.WorkflowExecutionID)
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

func (c *Controller) WaitForCapability(ctx context.Context, capability string, timeoutDuration time.Duration) error {
	// Create a context with timeout if not already set
	ctx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	c.lggr.Info().Msgf("Waiting for capability %s on all nodes...", capability)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out or context cancelled while waiting for capability %s: %w", capability, ctx.Err())
		case <-ticker.C:
			capInfos, err := c.List(ctx)
			if err != nil {
				c.lggr.Error().Err(err).Msgf("Failed to list capabilities while waiting for %s", capability)
				continue
			}

			allNodesHaveCapability := true
			for _, nodeInfo := range capInfos {
				hasCapability := false
				for _, singleCap := range nodeInfo.Capabilities {
					if singleCap.ID == capability {
						hasCapability = true
						break
					}
				}
				if !hasCapability {
					c.lggr.Debug().Msgf("Node %s does not have capability %s yet", nodeInfo.Node, capability)
					allNodesHaveCapability = false
					break
				}
			}

			if allNodesHaveCapability {
				c.lggr.Info().Msgf("All nodes now have capability %s", capability)
				return nil
			}
		}
	}
}

// GetTriggerSubscribers retrieves all subscribers for a specific trigger ID from all nodes
func (c *Controller) GetTriggerSubscribers(ctx context.Context, triggerID string) (map[string][]string, error) {
	subscribers := make(map[string][]string, 0)

	for _, client := range c.Nodes {
		resp, err := client.API.GetTriggerSubscribers(ctx, &pb2.GetTriggerSubscribersRequest{
			ID: triggerID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get trigger subscribers from node %s: %w", client.URL, err)
		}

		c.lggr.Debug().
			Str("node", client.URL).
			Str("triggerID", triggerID).
			Int("subscriberCount", len(resp.WorkflowIDs)).
			Msg("Retrieved trigger subscribers")

		subscribers[client.URL] = resp.WorkflowIDs
	}

	return subscribers, nil
}

// WaitForTriggerSubscribers waits until all nodes have at least one subscriber for the specified trigger
func (c *Controller) WaitForTriggerSubscribers(ctx context.Context, triggerID string, timeoutDuration time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	c.lggr.Info().Msgf("Waiting for subscribers on trigger %s for all nodes...", triggerID)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out or context cancelled while waiting for subscribers on trigger %s: %w", triggerID, ctx.Err())
		case <-ticker.C:
			subscribers, err := c.GetTriggerSubscribers(ctx, triggerID)
			if err != nil {
				c.lggr.Error().Err(err).Msgf("Failed to get trigger subscribers while waiting for %s", triggerID)
				continue
			}

			allNodesHaveSubscribers := true
			for nodeURL, workflowIDs := range subscribers {
				if len(workflowIDs) == 0 {
					c.lggr.Debug().Msgf("Node %s does not have subscribers for trigger %s yet", nodeURL, triggerID)
					allNodesHaveSubscribers = false
					break
				}
			}

			// Check if all nodes are represented in the subscribers map
			if len(subscribers) < len(c.Nodes) {
				missingNodes := []string{}
				for _, node := range c.Nodes {
					if _, exists := subscribers[node.URL]; !exists {
						missingNodes = append(missingNodes, node.URL)
					}
				}
				if len(missingNodes) > 0 {
					c.lggr.Debug().Msgf("Some nodes have no subscribers for trigger %s: %v", triggerID, missingNodes)
					allNodesHaveSubscribers = false
				}
			}

			if allNodesHaveSubscribers {
				c.lggr.Info().Msgf("All nodes now have subscribers for trigger %s", triggerID)
				return nil
			}
		}
	}
}

// HasCapability checks if a capability with the given ID exists on any node
func (c *Controller) HasCapability(ctx context.Context, capabilityID string) (bool, error) {
	capInfos, err := c.List(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to list capabilities: %w", err)
	}

	for _, nodeInfo := range capInfos {
		hasCapability := false
		for _, cap := range nodeInfo.Capabilities {
			if cap.ID == capabilityID {
				hasCapability = true
				break
			}
		}
		if !hasCapability {
			// If any node doesn't have the capability, return false
			return false, nil
		}
	}

	// All nodes have the capability
	return true, nil
}

func (c *Controller) DeleteCapability(ctx context.Context, capabilityID string) error {
	for _, client := range c.Nodes {
		_, err := client.API.RemoveCapability(ctx, &pb2.RemoveCapabilityRequest{
			ID: capabilityID,
		})
		if err != nil {
			return fmt.Errorf("failed to delete capability %s on %s: %w", capabilityID, client.URL, err)
		}
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

// StringToCapabilityType converts a string capability type to the corresponding integer constant
func StringToCapabilityType(typeStr string) pb2.CapabilityType {
	switch strings.ToUpper(typeStr) {
	case "TRIGGER":
		return pb2.CapabilityType_Trigger
	case "CONSENSUS":
		return pb2.CapabilityType_Consensus
	case "ACTION":
		return pb2.CapabilityType_Action
	case "TARGET":
		return pb2.CapabilityType_Target
	default:
		return pb2.CapabilityType_Unknown
	}
}
