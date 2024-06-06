package capability

import (
	"context"
	"errors"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type ActionCapabilityClient struct {
	*callbackExecutableClient
	*baseCapabilityClient
}

func NewActionCapabilityClient(brokerExt *net.BrokerExt, conn *grpc.ClientConn) capabilities.ActionCapability {
	return &ActionCapabilityClient{
		callbackExecutableClient: newCallbackExecutableClient(brokerExt, conn),
		baseCapabilityClient:     newBaseCapabilityClient(brokerExt, conn),
	}
}

type ConsensusCapabilityClient struct {
	*callbackExecutableClient
	*baseCapabilityClient
}

func NewConsensusCapabilityClient(brokerExt *net.BrokerExt, conn *grpc.ClientConn) capabilities.ConsensusCapability {
	return &ConsensusCapabilityClient{
		callbackExecutableClient: newCallbackExecutableClient(brokerExt, conn),
		baseCapabilityClient:     newBaseCapabilityClient(brokerExt, conn),
	}
}

type TargetCapabilityClient struct {
	*callbackExecutableClient
	*baseCapabilityClient
}

func NewTargetCapabilityClient(brokerExt *net.BrokerExt, conn *grpc.ClientConn) capabilities.TargetCapability {
	return &TargetCapabilityClient{
		callbackExecutableClient: newCallbackExecutableClient(brokerExt, conn),
		baseCapabilityClient:     newBaseCapabilityClient(brokerExt, conn),
	}
}

type TriggerCapabilityClient struct {
	*triggerExecutableClient
	*baseCapabilityClient
}

func NewTriggerCapabilityClient(brokerExt *net.BrokerExt, conn *grpc.ClientConn) capabilities.TriggerCapability {
	return &TriggerCapabilityClient{
		triggerExecutableClient: newTriggerExecutableClient(brokerExt, conn),
		baseCapabilityClient:    newBaseCapabilityClient(brokerExt, conn),
	}
}

type CallbackCapabilityClient struct {
	*callbackExecutableClient
	*baseCapabilityClient
}

type CallbackCapability interface {
	capabilities.CallbackExecutable
	capabilities.BaseCapability
}

func NewCallbackCapabilityClient(brokerExt *net.BrokerExt, conn *grpc.ClientConn) CallbackCapability {
	return &CallbackCapabilityClient{
		callbackExecutableClient: newCallbackExecutableClient(brokerExt, conn),
		baseCapabilityClient:     newBaseCapabilityClient(brokerExt, conn),
	}
}

func RegisterCallbackCapabilityServer(server *grpc.Server, broker net.Broker, brokerCfg net.BrokerConfig, impl CallbackCapability) error {
	bext := &net.BrokerExt{
		BrokerConfig: brokerCfg,
		Broker:       broker,
	}
	capabilitiespb.RegisterCallbackExecutableServer(server, newCallbackExecutableServer(bext, impl))
	capabilitiespb.RegisterBaseCapabilityServer(server, newBaseCapabilityServer(impl))
	return nil
}

func RegisterTriggerCapabilityServer(server *grpc.Server, broker net.Broker, brokerCfg net.BrokerConfig, impl capabilities.TriggerCapability) error {
	bext := &net.BrokerExt{
		BrokerConfig: brokerCfg,
		Broker:       broker,
	}
	capabilitiespb.RegisterTriggerExecutableServer(server, newTriggerExecutableServer(bext, impl))
	capabilitiespb.RegisterBaseCapabilityServer(server, newBaseCapabilityServer(impl))
	return nil
}

type baseCapabilityServer struct {
	capabilitiespb.UnimplementedBaseCapabilityServer

	impl capabilities.BaseCapability
}

func newBaseCapabilityServer(impl capabilities.BaseCapability) *baseCapabilityServer {
	return &baseCapabilityServer{impl: impl}
}

var _ capabilitiespb.BaseCapabilityServer = (*baseCapabilityServer)(nil)

func (c *baseCapabilityServer) Info(ctx context.Context, request *emptypb.Empty) (*capabilitiespb.CapabilityInfoReply, error) {
	info, err := c.impl.Info(ctx)
	if err != nil {
		return nil, err
	}

	return capabilityInfoToCapabilityInfoReply(info), nil
}

func capabilityInfoToCapabilityInfoReply(info capabilities.CapabilityInfo) *capabilitiespb.CapabilityInfoReply {
	var ct capabilitiespb.CapabilityType
	switch info.CapabilityType {
	case capabilities.CapabilityTypeTrigger:
		ct = capabilitiespb.CapabilityType_CAPABILITY_TYPE_TRIGGER
	case capabilities.CapabilityTypeAction:
		ct = capabilitiespb.CapabilityType_CAPABILITY_TYPE_ACTION
	case capabilities.CapabilityTypeConsensus:
		ct = capabilitiespb.CapabilityType_CAPABILITY_TYPE_CONSENSUS
	case capabilities.CapabilityTypeTarget:
		ct = capabilitiespb.CapabilityType_CAPABILITY_TYPE_TARGET
	}

	return &capabilitiespb.CapabilityInfoReply{
		Id:             info.ID,
		CapabilityType: ct,
		Description:    info.Description,
		Version:        info.Version,
	}
}

type baseCapabilityClient struct {
	grpc capabilitiespb.BaseCapabilityClient
	*net.BrokerExt
}

var _ capabilities.BaseCapability = (*baseCapabilityClient)(nil)

func newBaseCapabilityClient(brokerExt *net.BrokerExt, conn *grpc.ClientConn) *baseCapabilityClient {
	return &baseCapabilityClient{grpc: capabilitiespb.NewBaseCapabilityClient(conn), BrokerExt: brokerExt}
}

func (c *baseCapabilityClient) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	reply, err := c.grpc.Info(ctx, &emptypb.Empty{})
	if err != nil {
		return capabilities.CapabilityInfo{}, err
	}

	return capabilityInfoReplyToCapabilityInfo(reply)
}

func capabilityInfoReplyToCapabilityInfo(resp *capabilitiespb.CapabilityInfoReply) (capabilities.CapabilityInfo, error) {
	var ct capabilities.CapabilityType
	switch resp.CapabilityType {
	case capabilitiespb.CapabilityTypeTrigger:
		ct = capabilities.CapabilityTypeTrigger
	case capabilitiespb.CapabilityTypeAction:
		ct = capabilities.CapabilityTypeAction
	case capabilitiespb.CapabilityTypeConsensus:
		ct = capabilities.CapabilityTypeConsensus
	case capabilitiespb.CapabilityTypeTarget:
		ct = capabilities.CapabilityTypeTarget
	case capabilitiespb.CapabilityTypeUnknown:
		return capabilities.CapabilityInfo{}, fmt.Errorf("invalid capability type: %s", ct)
	}

	return capabilities.CapabilityInfo{
		ID:             resp.Id,
		CapabilityType: ct,
		Description:    resp.Description,
		Version:        resp.Version,
	}, nil
}

type triggerExecutableServer struct {
	capabilitiespb.UnimplementedTriggerExecutableServer
	*net.BrokerExt

	impl capabilities.TriggerExecutable
}

func newTriggerExecutableServer(brokerExt *net.BrokerExt, impl capabilities.TriggerExecutable) *triggerExecutableServer {
	return &triggerExecutableServer{
		impl:      impl,
		BrokerExt: brokerExt,
	}
}

var _ capabilitiespb.TriggerExecutableServer = (*triggerExecutableServer)(nil)

func (t *triggerExecutableServer) RegisterTrigger(request *capabilitiespb.CapabilityRequest,
	server capabilitiespb.TriggerExecutable_RegisterTriggerServer) error {
	req := pb.CapabilityRequestFromProto(request)
	responseCh, err := t.impl.RegisterTrigger(server.Context(), req)
	if err != nil {
		return fmt.Errorf("error registering trigger: %w", err)
	}

	defer func() {
		// Always attempt to unregister the trigger to ensure any related resources are cleaned up
		err = t.impl.UnregisterTrigger(server.Context(), req)
		if err != nil {
			t.Logger.Error("error unregistering trigger", "err", err)
		}
	}()

	for {
		select {
		case <-server.Context().Done():
			return nil
		case resp, ok := <-responseCh:
			if !ok {
				return nil
			}

			msg := &capabilitiespb.ResponseMessage{
				Message: &capabilitiespb.ResponseMessage_Response{
					Response: pb.CapabilityResponseToProto(resp),
				},
			}
			if err = server.Send(msg); err != nil {
				return fmt.Errorf("error sending response for trigger %s: %w", request, err)
			}
		}
	}
}

func (t *triggerExecutableServer) UnregisterTrigger(ctx context.Context, request *capabilitiespb.CapabilityRequest) (*emptypb.Empty, error) {
	if err := t.impl.UnregisterTrigger(ctx, pb.CapabilityRequestFromProto(request)); err != nil {
		return nil, fmt.Errorf("error unregistering trigger: %w", err)
	}

	return &emptypb.Empty{}, nil
}

type triggerExecutableClient struct {
	grpc capabilitiespb.TriggerExecutableClient
	*net.BrokerExt
}

func (t *triggerExecutableClient) RegisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	responseStream, err := t.grpc.RegisterTrigger(ctx, pb.CapabilityRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("error registering trigger: %w", err)
	}

	return forwardResponsesToChannel(ctx, t.Logger, req, responseStream.Recv)
}

func (t *triggerExecutableClient) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	_, err := t.grpc.UnregisterTrigger(ctx, pb.CapabilityRequestToProto(req))
	return err
}

func newTriggerExecutableClient(brokerExt *net.BrokerExt, conn *grpc.ClientConn) *triggerExecutableClient {
	return &triggerExecutableClient{grpc: capabilitiespb.NewTriggerExecutableClient(conn), BrokerExt: brokerExt}
}

type callbackExecutableServer struct {
	capabilitiespb.UnimplementedCallbackExecutableServer
	*net.BrokerExt

	impl capabilities.CallbackExecutable

	cancelFuncs map[string]func()
}

func newCallbackExecutableServer(brokerExt *net.BrokerExt, impl capabilities.CallbackExecutable) *callbackExecutableServer {
	return &callbackExecutableServer{
		impl:        impl,
		BrokerExt:   brokerExt,
		cancelFuncs: map[string]func(){},
	}
}

var _ capabilitiespb.CallbackExecutableServer = (*callbackExecutableServer)(nil)

func (c *callbackExecutableServer) RegisterToWorkflow(ctx context.Context, req *capabilitiespb.RegisterToWorkflowRequest) (*emptypb.Empty, error) {
	config := values.FromProto(req.Config)

	err := c.impl.RegisterToWorkflow(ctx, capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: req.Metadata.WorkflowId,
		},
		Config: config.(*values.Map),
	})
	return &emptypb.Empty{}, err
}

func (c *callbackExecutableServer) UnregisterFromWorkflow(ctx context.Context, req *capabilitiespb.UnregisterFromWorkflowRequest) (*emptypb.Empty, error) {
	config := values.FromProto(req.Config)

	err := c.impl.UnregisterFromWorkflow(ctx, capabilities.UnregisterFromWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: req.Metadata.WorkflowId,
		},
		Config: config.(*values.Map),
	})
	return &emptypb.Empty{}, err
}

func (c *callbackExecutableServer) Execute(req *capabilitiespb.CapabilityRequest, server capabilitiespb.CallbackExecutable_ExecuteServer) error {
	responseCh, err := c.impl.Execute(server.Context(), pb.CapabilityRequestFromProto(req))
	if err != nil {
		return fmt.Errorf("error executing capability request: %w", err)
	}

	err = server.Send(&capabilitiespb.ResponseMessage{
		Message: &capabilitiespb.ResponseMessage_Ack{
			Ack: &emptypb.Empty{},
		},
	})
	if err != nil {
		return fmt.Errorf("error sending ack: %w", err)
	}

	for resp := range responseCh {
		msg := &capabilitiespb.ResponseMessage{
			Message: &capabilitiespb.ResponseMessage_Response{
				Response: pb.CapabilityResponseToProto(resp),
			},
		}
		if err = server.Send(msg); err != nil {
			return fmt.Errorf("error sending response for execute request %s: %w", req, err)
		}
	}

	return nil
}

type callbackExecutableClient struct {
	grpc capabilitiespb.CallbackExecutableClient
	*net.BrokerExt
}

func newCallbackExecutableClient(brokerExt *net.BrokerExt, conn *grpc.ClientConn) *callbackExecutableClient {
	return &callbackExecutableClient{
		grpc:      capabilitiespb.NewCallbackExecutableClient(conn),
		BrokerExt: brokerExt,
	}
}

var _ capabilities.CallbackExecutable = (*callbackExecutableClient)(nil)

func (c *callbackExecutableClient) Execute(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	responseStream, err := c.grpc.Execute(ctx, pb.CapabilityRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("error executing capability request: %w", err)
	}

	resp, err := responseStream.Recv()
	if err != nil {
		return nil, fmt.Errorf("error waiting for ack: %w", err)
	}

	if _, ok := resp.GetMessage().(*capabilitiespb.ResponseMessage_Ack); !ok {
		return nil, fmt.Errorf("protocol error: first message received was not an ack: %+v", resp.GetMessage())
	}

	return forwardResponsesToChannel(ctx, c.Logger, req, responseStream.Recv)
}

func (c *callbackExecutableClient) UnregisterFromWorkflow(ctx context.Context, req capabilities.UnregisterFromWorkflowRequest) error {
	config := &values.Map{Underlying: map[string]values.Value{}}
	if req.Config != nil {
		config = req.Config
	}

	r := &capabilitiespb.UnregisterFromWorkflowRequest{
		Config: values.Proto(config),
		Metadata: &capabilitiespb.RegistrationMetadata{
			WorkflowId: req.Metadata.WorkflowID,
		},
	}

	_, err := c.grpc.UnregisterFromWorkflow(ctx, r)
	return err
}

func (c *callbackExecutableClient) RegisterToWorkflow(ctx context.Context, req capabilities.RegisterToWorkflowRequest) error {
	config := &values.Map{Underlying: map[string]values.Value{}}
	if req.Config != nil {
		config = req.Config
	}

	r := &capabilitiespb.RegisterToWorkflowRequest{
		Config: values.Proto(config),
		Metadata: &capabilitiespb.RegistrationMetadata{
			WorkflowId: req.Metadata.WorkflowID,
		},
	}

	_, err := c.grpc.RegisterToWorkflow(ctx, r)
	return err
}

func forwardResponsesToChannel(ctx context.Context, logger logger.Logger, req capabilities.CapabilityRequest, receive func() (*capabilitiespb.ResponseMessage, error)) (<-chan capabilities.CapabilityResponse, error) {
	responseCh := make(chan capabilities.CapabilityResponse)

	go func() {
		defer close(responseCh)
		for {
			message, err := receive()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				resp := capabilities.CapabilityResponse{
					Err: err,
				}
				select {
				case responseCh <- resp:
				case <-ctx.Done():
				}
				return
			}

			resp := message.GetResponse()
			if resp == nil {
				resp := capabilities.CapabilityResponse{
					Err: errors.New("unexpected message type when receiving response: expected response"),
				}
				select {
				case responseCh <- resp:
				case <-ctx.Done():
				}
				return
			}

			select {
			case responseCh <- pb.CapabilityResponseFromProto(resp):
			case <-ctx.Done():
				return
			}
		}
	}()

	return responseCh, nil
}
