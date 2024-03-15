package internal

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.CapabilitiesRegistry = (*capabilitiesRegistryClient)(nil)

type capabilitiesRegistryClient struct {
	*net.BrokerExt
	grpc pb.CapabilitiesRegistryClient
}

func (cr *capabilitiesRegistryClient) Get(ctx context.Context, ID string) (capabilities.BaseCapability, error) {
	req := &pb.GetRequest{
		Id: ID,
	}

	res, err := cr.grpc.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	conn, err := cr.Dial(res.CapabilityID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "Capability", ID: res.CapabilityID, Err: err}
	}
	client := newBaseCapabilityClient(cr.BrokerExt, conn)
	return client, nil
}

func (cr *capabilitiesRegistryClient) GetTrigger(ctx context.Context, ID string) (capabilities.TriggerCapability, error) {
	req := &pb.GetTriggerRequest{
		Id: ID,
	}

	res, err := cr.grpc.GetTrigger(ctx, req)
	if err != nil {
		return nil, err
	}

	conn, err := cr.Dial(res.CapabilityID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "GetTrigger", ID: res.CapabilityID, Err: err}
	}
	client := NewTriggerCapabilityClient(cr.BrokerExt, conn)
	return client, nil
}

func (cr *capabilitiesRegistryClient) GetAction(ctx context.Context, ID string) (capabilities.ActionCapability, error) {
	req := &pb.GetActionRequest{
		Id: ID,
	}

	res, err := cr.grpc.GetAction(ctx, req)
	if err != nil {
		return nil, err
	}
	conn, err := cr.Dial(res.CapabilityID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "GetAction", ID: res.CapabilityID, Err: err}
	}
	client := NewActionCapabilityClient(cr.BrokerExt, conn)
	return client, nil
}

func (cr *capabilitiesRegistryClient) GetConsensus(ctx context.Context, ID string) (capabilities.ConsensusCapability, error) {
	req := &pb.GetConsensusRequest{
		Id: ID,
	}

	res, err := cr.grpc.GetConsensus(ctx, req)
	if err != nil {
		return nil, err
	}

	conn, err := cr.Dial(res.CapabilityID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "GetConsensus", ID: res.CapabilityID, Err: err}
	}
	client := NewConsensusCapabilityClient(cr.BrokerExt, conn)
	return client, nil
}

func (cr *capabilitiesRegistryClient) GetTarget(ctx context.Context, ID string) (capabilities.TargetCapability, error) {
	req := &pb.GetTargetRequest{
		Id: ID,
	}

	res, err := cr.grpc.GetTarget(ctx, req)
	if err != nil {
		return nil, err
	}

	conn, err := cr.Dial(res.CapabilityID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "GetTarget", ID: res.CapabilityID, Err: err}
	}
	client := NewTargetCapabilityClient(cr.BrokerExt, conn)
	return client, nil
}

func (cr *capabilitiesRegistryClient) List(ctx context.Context) ([]capabilities.BaseCapability, error) {
	res, err := cr.grpc.List(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	var clients []capabilities.BaseCapability
	for _, id := range res.CapabilityID {
		conn, err := cr.Dial(id)
		if err != nil {
			return nil, net.ErrConnDial{Name: "List", ID: id, Err: err}
		}
		client := newBaseCapabilityClient(cr.BrokerExt, conn)
		clients = append(clients, client)
	}

	return clients, nil
}

func (cr *capabilitiesRegistryClient) Add(ctx context.Context, c capabilities.BaseCapability) error {
	info, err := c.Info(ctx)
	if err != nil {
		return err
	}

	// Check the capability and the CapabilityType match here as the ServeNew method does not return an error
	err = validateCapability(c, info.CapabilityType)
	if err != nil {
		return err
	}

	var cRes net.Resource
	id, cRes, err := cr.ServeNew(info.ID, func(s *grpc.Server) {
		pbRegisterCapability(s, cr.BrokerExt, c, info.CapabilityType)
	})
	if err != nil {
		return err
	}

	_, err = cr.grpc.Add(ctx, &pb.AddRequest{
		CapabilityID: id,
		Type:         pb.ExecuteAPIType(getExecuteAPIType(info.CapabilityType)),
	})
	if err != nil {
		cRes.Close()
		return err
	}
	return nil
}

func NewCapabilitiesRegistryClient(cc grpc.ClientConnInterface, b *net.BrokerExt) *capabilitiesRegistryClient {
	return &capabilitiesRegistryClient{grpc: pb.NewCapabilitiesRegistryClient(cc), BrokerExt: b.WithName("CapabilitiesRegistryClient")}
}

var _ pb.CapabilitiesRegistryServer = (*capabilitiesRegistryServer)(nil)

type capabilitiesRegistryServer struct {
	pb.UnimplementedCapabilitiesRegistryServer
	*net.BrokerExt
	impl types.CapabilitiesRegistry
}

func (c *capabilitiesRegistryServer) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetReply, error) {
	capability, err := c.impl.Get(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	info, err := capability.Info(ctx)
	if err != nil {
		return nil, err
	}

	// Check the capability and the CapabilityType match here as the ServeNew method does not return an error
	err = validateCapability(capability, info.CapabilityType)
	if err != nil {
		return nil, err
	}

	id, _, err := c.ServeNew("Get", func(s *grpc.Server) {
		pbRegisterCapability(s, c.BrokerExt, capability, info.CapabilityType)
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetReply{
		CapabilityID: id,
		Type:         pb.ExecuteAPIType(getExecuteAPIType(info.CapabilityType)),
	}, nil
}

func (c *capabilitiesRegistryServer) GetTrigger(ctx context.Context, request *pb.GetTriggerRequest) (*pb.GetTriggerReply, error) {
	capability, err := c.impl.GetTrigger(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	// Check the capability and the CapabilityType match here as the ServeNew method does not return an error
	err = validateCapability(capability, capabilities.CapabilityTypeTrigger)
	if err != nil {
		return nil, err
	}

	id, _, err := c.ServeNew("GetTrigger", func(s *grpc.Server) {
		pbRegisterCapability(s, c.BrokerExt, capability, capabilities.CapabilityTypeTrigger)
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetTriggerReply{
		CapabilityID: id,
	}, nil
}

func (c *capabilitiesRegistryServer) GetAction(ctx context.Context, request *pb.GetActionRequest) (*pb.GetActionReply, error) {
	capability, err := c.impl.GetAction(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	// Check the capability and the CapabilityType match here as the ServeNew method does not return an error
	err = validateCapability(capability, capabilities.CapabilityTypeAction)
	if err != nil {
		return nil, err
	}

	id, _, err := c.ServeNew("GetAction", func(s *grpc.Server) {
		pbRegisterCapability(s, c.BrokerExt, capability, capabilities.CapabilityTypeAction)
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetActionReply{
		CapabilityID: id,
	}, nil
}

func (c *capabilitiesRegistryServer) GetConsensus(ctx context.Context, request *pb.GetConsensusRequest) (*pb.GetConsensusReply, error) {
	capability, err := c.impl.GetConsensus(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	// Check the capability and the CapabilityType match here as the ServeNew method does not return an error
	err = validateCapability(capability, capabilities.CapabilityTypeConsensus)
	if err != nil {
		return nil, err
	}

	id, _, err := c.ServeNew("GetConsensus", func(s *grpc.Server) {
		pbRegisterCapability(s, c.BrokerExt, capability, capabilities.CapabilityTypeConsensus)
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetConsensusReply{
		CapabilityID: id,
	}, nil
}

func (c *capabilitiesRegistryServer) GetTarget(ctx context.Context, request *pb.GetTargetRequest) (*pb.GetTargetReply, error) {
	capability, err := c.impl.GetTarget(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	// Check the capability and the CapabilityType match here as the ServeNew method does not return an error
	err = validateCapability(capability, capabilities.CapabilityTypeTarget)
	if err != nil {
		return nil, err
	}

	id, _, err := c.ServeNew("GetTarget", func(s *grpc.Server) {
		pbRegisterCapability(s, c.BrokerExt, capability, capabilities.CapabilityTypeTarget)
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetTargetReply{
		CapabilityID: id,
	}, nil
}

func (c *capabilitiesRegistryServer) List(ctx context.Context, _ *emptypb.Empty) (*pb.ListReply, error) {
	capabilities, err := c.impl.List(ctx)
	if err != nil {
		return nil, err
	}

	reply := &pb.ListReply{}

	var resources []net.Resource
	for _, cap := range capabilities {
		info, err := cap.Info(ctx)
		if err != nil {
			c.CloseAll(resources...)
			return nil, err
		}

		// Check the capability and the CapabilityType match here as the ServeNew method does not return an error
		err = validateCapability(cap, info.CapabilityType)
		if err != nil {
			c.CloseAll(resources...)
			return nil, err
		}

		id, res, err := c.ServeNew("List", func(s *grpc.Server) {
			pbRegisterCapability(s, c.BrokerExt, cap, info.CapabilityType)
		})
		if err != nil {
			c.CloseAll(resources...)
			return nil, err
		}
		resources = append(resources, res)
		reply.CapabilityID = append(reply.CapabilityID, id)
	}

	return reply, nil
}

func (c *capabilitiesRegistryServer) Add(ctx context.Context, request *pb.AddRequest) (*emptypb.Empty, error) {
	conn, err := c.Dial(request.CapabilityID)
	if err != nil {
		return &emptypb.Empty{}, net.ErrConnDial{Name: "Add", ID: request.CapabilityID, Err: err}
	}
	var client capabilities.BaseCapability

	switch request.Type {
	case pb.ExecuteAPIType_EXECUTE_API_TYPE_TRIGGER:
		client = NewTriggerCapabilityClient(c.BrokerExt, conn)
	case pb.ExecuteAPIType_EXECUTE_API_TYPE_CALLBACK:
		client = NewCallbackCapabilityClient(c.BrokerExt, conn)
	default:
		return nil, fmt.Errorf("unknown execute type %d", request.Type)
	}

	err = c.impl.Add(ctx, client)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func NewCapabilitiesRegistryServer(b *net.BrokerExt, i types.CapabilitiesRegistry) *capabilitiesRegistryServer {
	return &capabilitiesRegistryServer{
		BrokerExt: b.WithName("CapabilitiesRegistryServer"),
		impl:      i,
	}
}

func validateCapability(impl capabilities.BaseCapability, t capabilities.CapabilityType) error {
	switch t {
	case capabilities.CapabilityTypeTrigger:
		_, ok := impl.(capabilities.TriggerCapability)
		if !ok {
			return fmt.Errorf("expected TriggerCapability but got %T", impl)
		}
	case capabilities.CapabilityTypeAction:
		_, ok := impl.(capabilities.ActionCapability)
		if !ok {
			return fmt.Errorf("expected ActionCapability but got %T", impl)
		}
	case capabilities.CapabilityTypeConsensus:
		_, ok := impl.(capabilities.ConsensusCapability)
		if !ok {
			return fmt.Errorf("expected ConsensusCapability but got %T", impl)
		}
	case capabilities.CapabilityTypeTarget:
		_, ok := impl.(capabilities.TargetCapability)
		if !ok {
			return fmt.Errorf("expected TargetCapability but got %T", impl)
		}
	}
	return nil
}

// pbRegisterCapability registers the server with the correct capability based on capability type, this method assumes
// that the capability has already been validated with validateCapability.
func pbRegisterCapability(s *grpc.Server, b *net.BrokerExt, impl capabilities.BaseCapability, t capabilities.CapabilityType) {
	switch t {
	case capabilities.CapabilityTypeTrigger:
		i, _ := impl.(capabilities.TriggerCapability)
		pb.RegisterTriggerExecutableServer(s, &triggerExecutableServer{
			BrokerExt:   b,
			impl:        i,
			cancelFuncs: map[string]func(){},
		})
	case capabilities.CapabilityTypeAction:
		i, _ := impl.(capabilities.ActionCapability)

		pb.RegisterCallbackExecutableServer(s, &callbackExecutableServer{
			BrokerExt:   b,
			impl:        i,
			cancelFuncs: map[string]func(){},
		})
	case capabilities.CapabilityTypeConsensus:
		i, _ := impl.(capabilities.ConsensusCapability)

		pb.RegisterCallbackExecutableServer(s, &callbackExecutableServer{
			BrokerExt:   b,
			impl:        i,
			cancelFuncs: map[string]func(){},
		})
	case capabilities.CapabilityTypeTarget:
		i, _ := impl.(capabilities.TargetCapability)
		pb.RegisterCallbackExecutableServer(s, &callbackExecutableServer{
			BrokerExt:   b,
			impl:        i,
			cancelFuncs: map[string]func(){},
		})
	}
	pb.RegisterBaseCapabilityServer(s, newBaseCapabilityServer(impl))
}

func getExecuteAPIType(c capabilities.CapabilityType) int32 {
	switch c {
	case capabilities.CapabilityTypeTrigger:
		return 1
	case capabilities.CapabilityTypeAction:
		return 2
	case capabilities.CapabilityTypeConsensus:
		return 2
	case capabilities.CapabilityTypeTarget:
		return 2
	default:
		return 0
	}
}
