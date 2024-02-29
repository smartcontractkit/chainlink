package internal

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/grpc"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type ActionCapabilityClient struct {
	*callbackExecutableClient
	*baseCapabilityClient
}

func NewActionCapabilityClient(brokerExt *BrokerExt, conn *grpc.ClientConn) capabilities.ActionCapability {
	return &ActionCapabilityClient{
		callbackExecutableClient: newCallbackExecutableClient(brokerExt, conn),
		baseCapabilityClient:     newBaseCapabilityClient(brokerExt, conn),
	}
}

type ConsensusCapabilityClient struct {
	*callbackExecutableClient
	*baseCapabilityClient
}

func NewConsensusCapabilityClient(brokerExt *BrokerExt, conn *grpc.ClientConn) capabilities.ConsensusCapability {
	return &ConsensusCapabilityClient{
		callbackExecutableClient: newCallbackExecutableClient(brokerExt, conn),
		baseCapabilityClient:     newBaseCapabilityClient(brokerExt, conn),
	}
}

type TargetCapabilityClient struct {
	*callbackExecutableClient
	*baseCapabilityClient
}

func NewTargetCapabilityClient(brokerExt *BrokerExt, conn *grpc.ClientConn) capabilities.TargetCapability {
	return &TargetCapabilityClient{
		callbackExecutableClient: newCallbackExecutableClient(brokerExt, conn),
		baseCapabilityClient:     newBaseCapabilityClient(brokerExt, conn),
	}
}

type TriggerCapabilityClient struct {
	*triggerExecutableClient
	*baseCapabilityClient
}

func NewTriggerCapabilityClient(brokerExt *BrokerExt, conn *grpc.ClientConn) capabilities.TriggerCapability {
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

func NewCallbackCapabilityClient(brokerExt *BrokerExt, conn *grpc.ClientConn) CallbackCapability {
	return &CallbackCapabilityClient{
		callbackExecutableClient: newCallbackExecutableClient(brokerExt, conn),
		baseCapabilityClient:     newBaseCapabilityClient(brokerExt, conn),
	}
}

func RegisterCallbackCapabilityServer(server *grpc.Server, broker Broker, brokerCfg BrokerConfig, impl CallbackCapability) error {
	bext := &BrokerExt{
		BrokerConfig: brokerCfg,
		Broker:       broker,
	}
	pb.RegisterCallbackExecutableServer(server, newCallbackExecutableServer(bext, impl))
	pb.RegisterBaseCapabilityServer(server, newBaseCapabilityServer(impl))
	return nil
}

func RegisterTriggerCapabilityServer(server *grpc.Server, broker Broker, brokerCfg BrokerConfig, impl capabilities.TriggerCapability) error {
	bext := &BrokerExt{
		BrokerConfig: brokerCfg,
		Broker:       broker,
	}
	pb.RegisterTriggerExecutableServer(server, newTriggerExecutableServer(bext, impl))
	pb.RegisterBaseCapabilityServer(server, newBaseCapabilityServer(impl))
	return nil
}

type baseCapabilityServer struct {
	pb.UnimplementedBaseCapabilityServer

	impl capabilities.BaseCapability
}

func newBaseCapabilityServer(impl capabilities.BaseCapability) *baseCapabilityServer {
	return &baseCapabilityServer{impl: impl}
}

var _ pb.BaseCapabilityServer = (*baseCapabilityServer)(nil)

func (c *baseCapabilityServer) Info(ctx context.Context, request *emptypb.Empty) (*pb.CapabilityInfoReply, error) {
	info, err := c.impl.Info(ctx)
	if err != nil {
		return nil, err
	}

	var ct pb.CapabilityType
	switch info.CapabilityType {
	case capabilities.CapabilityTypeTrigger:
		ct = pb.CapabilityType_CAPABILITY_TYPE_TRIGGER
	case capabilities.CapabilityTypeAction:
		ct = pb.CapabilityType_CAPABILITY_TYPE_ACTION
	case capabilities.CapabilityTypeConsensus:
		ct = pb.CapabilityType_CAPABILITY_TYPE_CONSENSUS
	case capabilities.CapabilityTypeTarget:
		ct = pb.CapabilityType_CAPABILITY_TYPE_TARGET
	}

	return &pb.CapabilityInfoReply{
		Id:             info.ID,
		CapabilityType: ct,
		Description:    info.Description,
		Version:        info.Version,
	}, nil
}

type baseCapabilityClient struct {
	grpc pb.BaseCapabilityClient
	*BrokerExt
}

var _ capabilities.BaseCapability = (*baseCapabilityClient)(nil)

func newBaseCapabilityClient(brokerExt *BrokerExt, conn *grpc.ClientConn) *baseCapabilityClient {
	return &baseCapabilityClient{grpc: pb.NewBaseCapabilityClient(conn), BrokerExt: brokerExt}
}

func (c *baseCapabilityClient) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	resp, err := c.grpc.Info(ctx, &emptypb.Empty{})
	if err != nil {
		return capabilities.CapabilityInfo{}, err
	}

	var ct capabilities.CapabilityType
	switch resp.CapabilityType {
	case pb.CapabilityTypeTrigger:
		ct = capabilities.CapabilityTypeTrigger
	case pb.CapabilityTypeAction:
		ct = capabilities.CapabilityTypeAction
	case pb.CapabilityTypeConsensus:
		ct = capabilities.CapabilityTypeConsensus
	case pb.CapabilityTypeTarget:
		ct = capabilities.CapabilityTypeTarget
	case pb.CapabilityTypeUnknown:
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
	pb.UnimplementedTriggerExecutableServer
	*BrokerExt

	impl capabilities.TriggerExecutable

	cancelFuncs map[string]func()
}

func newTriggerExecutableServer(brokerExt *BrokerExt, impl capabilities.TriggerExecutable) *triggerExecutableServer {
	return &triggerExecutableServer{
		impl:        impl,
		BrokerExt:   brokerExt,
		cancelFuncs: map[string]func(){},
	}
}

var _ pb.TriggerExecutableServer = (*triggerExecutableServer)(nil)

func (t *triggerExecutableServer) RegisterTrigger(ctx context.Context, request *pb.RegisterTriggerRequest) (*emptypb.Empty, error) {
	ch := make(chan capabilities.CapabilityResponse)

	conn, err := t.Dial(request.CallbackId)
	if err != nil {
		return nil, err
	}

	connCtx, connCancel := context.WithCancel(context.Background())
	go callbackIssuer(connCtx, pb.NewCallbackClient(conn), ch, t.Logger)

	cr := request.CapabilityRequest
	md := cr.Metadata

	config, err := values.FromProto(cr.Config)
	if err != nil {
		connCancel()
		return nil, err
	}

	inputs, err := values.FromProto(cr.Inputs)
	if err != nil {
		connCancel()
		return nil, err
	}

	req := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          md.WorkflowId,
			WorkflowExecutionID: md.WorkflowExecutionId,
		},
		Config: config.(*values.Map),
		Inputs: inputs.(*values.Map),
	}

	err = t.impl.RegisterTrigger(ctx, ch, req)
	if err != nil {
		connCancel()
		return nil, err
	}

	t.cancelFuncs[md.WorkflowId] = connCancel
	return &emptypb.Empty{}, nil
}

func (t *triggerExecutableServer) UnregisterTrigger(ctx context.Context, request *pb.UnregisterTriggerRequest) (*emptypb.Empty, error) {
	req := request.CapabilityRequest
	md := req.Metadata

	config, err := values.FromProto(req.Config)
	if err != nil {
		return nil, err
	}

	inputs, err := values.FromProto(req.Inputs)
	if err != nil {
		return nil, err
	}

	err = t.impl.UnregisterTrigger(ctx, capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          md.WorkflowId,
			WorkflowExecutionID: md.WorkflowExecutionId,
		},
		Inputs: inputs.(*values.Map),
		Config: config.(*values.Map),
	})
	if err != nil {
		return nil, err
	}

	cancelFunc := t.cancelFuncs[md.WorkflowId]
	if cancelFunc != nil {
		cancelFunc()
	}

	return &emptypb.Empty{}, nil
}

type triggerExecutableClient struct {
	grpc pb.TriggerExecutableClient
	*BrokerExt
}

var _ capabilities.TriggerExecutable = (*triggerExecutableClient)(nil)

func (t *triggerExecutableClient) RegisterTrigger(ctx context.Context, callback chan<- capabilities.CapabilityResponse, req capabilities.CapabilityRequest) error {
	cid, res, err := t.ServeNew("Callback", func(s *grpc.Server) {
		pb.RegisterCallbackServer(s, newCallbackServer(callback))
	})
	if err != nil {
		return err
	}

	reqPb, err := toProto(req)
	if err != nil {
		t.CloseAll(res)
		return err
	}

	r := &pb.RegisterTriggerRequest{
		CallbackId:        cid,
		CapabilityRequest: reqPb,
	}

	_, err = t.grpc.RegisterTrigger(ctx, r)
	if err != nil {
		t.CloseAll(res)
	}
	return err
}

func (t *triggerExecutableClient) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	reqPb, err := toProto(req)
	if err != nil {
		return err
	}

	r := &pb.UnregisterTriggerRequest{
		CapabilityRequest: reqPb,
	}

	_, err = t.grpc.UnregisterTrigger(ctx, r)
	return err
}

func newTriggerExecutableClient(brokerExt *BrokerExt, conn *grpc.ClientConn) *triggerExecutableClient {
	return &triggerExecutableClient{grpc: pb.NewTriggerExecutableClient(conn), BrokerExt: brokerExt}
}

type callbackExecutableServer struct {
	pb.UnimplementedCallbackExecutableServer
	*BrokerExt

	impl capabilities.CallbackExecutable

	cancelFuncs map[string]func()
}

func newCallbackExecutableServer(brokerExt *BrokerExt, impl capabilities.CallbackExecutable) *callbackExecutableServer {
	return &callbackExecutableServer{
		impl:        impl,
		BrokerExt:   brokerExt,
		cancelFuncs: map[string]func(){},
	}
}

var _ pb.CallbackExecutableServer = (*callbackExecutableServer)(nil)

func (c *callbackExecutableServer) RegisterToWorkflow(ctx context.Context, req *pb.RegisterToWorkflowRequest) (*emptypb.Empty, error) {
	config, err := values.FromProto(req.Config)
	if err != nil {
		return nil, err
	}

	err = c.impl.RegisterToWorkflow(ctx, capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: req.Metadata.WorkflowId,
		},
		Config: config.(*values.Map),
	})
	return &emptypb.Empty{}, err
}

func (c *callbackExecutableServer) UnregisterFromWorkflow(ctx context.Context, req *pb.UnregisterFromWorkflowRequest) (*emptypb.Empty, error) {
	config, err := values.FromProto(req.Config)
	if err != nil {
		return nil, err
	}

	err = c.impl.UnregisterFromWorkflow(ctx, capabilities.UnregisterFromWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: req.Metadata.WorkflowId,
		},
		Config: config.(*values.Map),
	})
	return &emptypb.Empty{}, err
}

func (c *callbackExecutableServer) Execute(ctx context.Context, req *pb.ExecuteRequest) (*emptypb.Empty, error) {
	ch := make(chan capabilities.CapabilityResponse)

	conn, err := c.Dial(req.CallbackId)
	if err != nil {
		return nil, err
	}

	connCtx, connCancel := context.WithCancel(context.Background())
	go callbackIssuer(connCtx, pb.NewCallbackClient(conn), ch, c.Logger)

	cr := req.CapabilityRequest
	md := cr.Metadata

	config, err := values.FromProto(cr.Config)
	if err != nil {
		connCancel()
		return nil, err
	}

	inputs, err := values.FromProto(cr.Inputs)
	if err != nil {
		connCancel()
		return nil, err
	}

	r := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          md.WorkflowId,
			WorkflowExecutionID: md.WorkflowExecutionId,
		},
		Config: config.(*values.Map),
		Inputs: inputs.(*values.Map),
	}

	err = c.impl.Execute(ctx, ch, r)
	if err != nil {
		connCancel()
		return nil, err
	}

	c.cancelFuncs[md.WorkflowId] = connCancel
	return &emptypb.Empty{}, nil
}

type callbackExecutableClient struct {
	grpc pb.CallbackExecutableClient
	*BrokerExt
}

func newCallbackExecutableClient(brokerExt *BrokerExt, conn *grpc.ClientConn) *callbackExecutableClient {
	return &callbackExecutableClient{
		grpc:      pb.NewCallbackExecutableClient(conn),
		BrokerExt: brokerExt,
	}
}

var _ capabilities.CallbackExecutable = (*callbackExecutableClient)(nil)

func toProto(req capabilities.CapabilityRequest) (*pb.CapabilityRequest, error) {
	inputs := &values.Map{Underlying: map[string]values.Value{}}
	if req.Inputs != nil {
		inputs = req.Inputs
	}

	config := &values.Map{Underlying: map[string]values.Value{}}
	if req.Config != nil {
		config = req.Config
	}

	return &pb.CapabilityRequest{
		Metadata: &pb.RequestMetadata{
			WorkflowId:          req.Metadata.WorkflowID,
			WorkflowExecutionId: req.Metadata.WorkflowExecutionID,
		},
		Inputs: values.Proto(inputs),
		Config: values.Proto(config),
	}, nil
}

func (c *callbackExecutableClient) Execute(ctx context.Context, callback chan<- capabilities.CapabilityResponse, req capabilities.CapabilityRequest) error {
	cid, res, err := c.ServeNew("Callback", func(s *grpc.Server) {
		pb.RegisterCallbackServer(s, newCallbackServer(callback))
	})
	if err != nil {
		return err
	}

	reqPb, err := toProto(req)
	if err != nil {
		c.CloseAll(res)
		return nil
	}

	r := &pb.ExecuteRequest{
		CallbackId:        cid,
		CapabilityRequest: reqPb,
	}

	_, err = c.grpc.Execute(ctx, r)
	if err != nil {
		c.CloseAll(res)
	}
	return err
}

func (c *callbackExecutableClient) UnregisterFromWorkflow(ctx context.Context, req capabilities.UnregisterFromWorkflowRequest) error {
	config := &values.Map{Underlying: map[string]values.Value{}}
	if req.Config != nil {
		config = req.Config
	}

	r := &pb.UnregisterFromWorkflowRequest{
		Config: values.Proto(config),
		Metadata: &pb.RegistrationMetadata{
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

	r := &pb.RegisterToWorkflowRequest{
		Config: values.Proto(config),
		Metadata: &pb.RegistrationMetadata{
			WorkflowId: req.Metadata.WorkflowID,
		},
	}

	_, err := c.grpc.RegisterToWorkflow(ctx, r)
	return err
}

type callbackServer struct {
	pb.UnimplementedCallbackServer
	ch chan<- capabilities.CapabilityResponse

	isClosed bool
	mu       sync.RWMutex
}

func newCallbackServer(ch chan<- capabilities.CapabilityResponse) *callbackServer {
	return &callbackServer{ch: ch}
}

func (c *callbackServer) SendResponse(ctx context.Context, req *pb.CapabilityResponse) (*emptypb.Empty, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.isClosed {
		return nil, errors.New("cannot send response: the underlying channel has been closed")
	}

	val, err := values.FromProto(req.Value)
	if err != nil {
		return nil, err
	}

	err = nil
	if req.Error != "" {
		err = errors.New(req.Error)
	}
	resp := capabilities.CapabilityResponse{
		Value: val,
		Err:   err,
	}
	c.ch <- resp
	return &emptypb.Empty{}, nil
}

func (c *callbackServer) CloseCallback(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	close(c.ch)
	c.isClosed = true
	return &emptypb.Empty{}, nil
}

func callbackIssuer(ctx context.Context, client pb.CallbackClient, callbackChannel chan capabilities.CapabilityResponse, logger logger.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp, isOpen := <-callbackChannel:
			if !isOpen {
				_, err := client.CloseCallback(ctx, &emptypb.Empty{})
				if err != nil {
					logger.Error("could not close upstream callback", err)
				}
				return
			}

			errStr := ""
			if resp.Err != nil {
				errStr = resp.Err.Error()
			}

			cr := &pb.CapabilityResponse{
				Error: errStr,
				Value: values.Proto(resp.Value),
			}

			_, err := client.SendResponse(ctx, cr)
			if err != nil {
				logger.Error("error sending callback response", err)
			}
		}
	}
}
